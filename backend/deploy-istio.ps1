# PowerShell script to deploy LDAP Manager Service with Istio

Write-Host "🚀 Deploying LDAP Manager Service with Istio..." -ForegroundColor Cyan
Write-Host ""

# Check if kubectl is available
Write-Host "📋 Checking prerequisites..." -ForegroundColor Yellow
try {
    kubectl version --client | Out-Null
} catch {
    Write-Host "❌ kubectl is not installed or not in PATH" -ForegroundColor Red
    exit 1
}

# Check if Istio is installed
try {
    kubectl get namespace istio-system 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Istio is not installed. Please install Istio first." -ForegroundColor Red
        Write-Host "   Run: istioctl install --set profile=default -y" -ForegroundColor Yellow
        exit 1
    }
    Write-Host "✅ Istio is installed" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to check Istio installation" -ForegroundColor Red
    exit 1
}

# Check if namespace exists, create if not
Write-Host ""
try {
    kubectl get namespace dev-platform 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Host "📦 Creating dev-platform namespace..." -ForegroundColor Yellow
        kubectl create namespace dev-platform
    }
} catch {
    Write-Host "📦 Creating dev-platform namespace..." -ForegroundColor Yellow
    kubectl create namespace dev-platform
}

# Enable Istio injection
Write-Host "🔧 Enabling Istio sidecar injection for dev-platform namespace..." -ForegroundColor Yellow
kubectl label namespace dev-platform istio-injection=enabled --overwrite
Write-Host "✅ Istio injection enabled" -ForegroundColor Green

# Deploy OpenLDAP
Write-Host ""
Write-Host "🗄️ Deploying OpenLDAP..." -ForegroundColor Yellow
kubectl apply -f k8s/01-openldap.yaml
Write-Host "✅ OpenLDAP deployed" -ForegroundColor Green

Write-Host ""
Write-Host "⏳ Waiting for OpenLDAP to be ready..." -ForegroundColor Yellow
kubectl wait --for=condition=ready pod -l app=openldap1 -n dev-platform --timeout=180s
if ($LASTEXITCODE -ne 0) {
    Write-Host "⚠️  OpenLDAP not ready within timeout, continuing anyway..." -ForegroundColor Yellow
} else {
    Write-Host "✅ OpenLDAP is ready" -ForegroundColor Green
}

# Run LDAP initialization job
Write-Host ""
Write-Host "🔧 Running LDAP initialization job..." -ForegroundColor Yellow
kubectl apply -f k8s/02-ldap-init.yaml
Write-Host ""
Write-Host "⏳ Waiting for LDAP init job to complete..." -ForegroundColor Yellow
kubectl wait --for=condition=complete job/ldap-init-structure -n dev-platform --timeout=180s
if ($LASTEXITCODE -ne 0) {
    Write-Host "⚠️  LDAP init job did not complete within timeout" -ForegroundColor Yellow
} else {
    Write-Host "✅ LDAP initialization completed" -ForegroundColor Green
}

# Deploy LDAP Manager
Write-Host ""
Write-Host "📝 Deploying ConfigMap and Secret..." -ForegroundColor Yellow
kubectl apply -f k8s/03-ldap-manager.yaml
Write-Host "✅ ConfigMap and Secret deployed" -ForegroundColor Green

# Wait for pods to be ready
Write-Host ""
Write-Host "⏳ Waiting for LDAP Manager pods to be ready (with Istio sidecar)..." -ForegroundColor Yellow
kubectl wait --for=condition=ready pod -l app=ldap-manager -n dev-platform --timeout=120s
if ($LASTEXITCODE -ne 0) {
    Write-Host "⚠️  Pods not ready within timeout, continuing anyway..." -ForegroundColor Yellow
} else {
    Write-Host "✅ LDAP Manager pods are ready" -ForegroundColor Green
}

# Deploy Istio configuration
Write-Host ""
Write-Host "🌐 Deploying Istio Gateway, VirtualService, and Policies..." -ForegroundColor Yellow
kubectl apply -f k8s/04-istio-config.yaml
Write-Host "✅ Istio configuration deployed" -ForegroundColor Green

# Verify deployment
Write-Host ""
Write-Host "🔍 Verifying deployment..." -ForegroundColor Cyan
Write-Host ""

# Check pods
Write-Host "Pods:" -ForegroundColor Yellow
kubectl get pods -n dev-platform -l app=ldap-manager

Write-Host ""
Write-Host "Gateway:" -ForegroundColor Yellow
kubectl get gateway -n dev-platform

Write-Host ""
Write-Host "VirtualService:" -ForegroundColor Yellow
kubectl get virtualservice -n dev-platform

Write-Host ""
Write-Host "DestinationRule:" -ForegroundColor Yellow
kubectl get destinationrule -n dev-platform

Write-Host ""
Write-Host "PeerAuthentication:" -ForegroundColor Yellow
kubectl get peerauthentication -n dev-platform

# Get Istio Ingress Gateway info
Write-Host ""
Write-Host "🌍 Istio Ingress Gateway Info:" -ForegroundColor Cyan

try {
    $INGRESS_HOST = kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}' 2>$null
    if ([string]::IsNullOrEmpty($INGRESS_HOST)) {
        $INGRESS_HOST = "pending"
    }
} catch {
    $INGRESS_HOST = "pending"
}

$INGRESS_PORT = kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}'

Write-Host "   Ingress Host: $INGRESS_HOST" -ForegroundColor White
Write-Host "   Ingress Port: $INGRESS_PORT" -ForegroundColor White

if ($INGRESS_HOST -ne "pending") {
    Write-Host ""
    Write-Host "🧪 Testing endpoints..." -ForegroundColor Cyan
    Write-Host ""

    # Test health endpoint
    Write-Host -NoNewline "Testing /health... "
    try {
        $response = Invoke-WebRequest -Uri "http://$INGRESS_HOST:$INGRESS_PORT/health" -Headers @{"Host"="ldap-manager.localhost"} -UseBasicParsing -TimeoutSec 5 -ErrorAction Stop
        if ($response.StatusCode -eq 200) {
            Write-Host "✅ OK" -ForegroundColor Green
        } else {
            Write-Host "⚠️  Not ready yet (Status: $($response.StatusCode))" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "⚠️  Not ready yet" -ForegroundColor Yellow
    }

    Write-Host ""
    Write-Host -NoNewline "Testing /graphql (login mutation)... "
    try {
        $body = '{"query":"mutation { login(uid: \"john.doe\", password: \"password123\") { token user { uid mail } } }"}'
        $response = Invoke-WebRequest -Uri "http://$INGRESS_HOST:$INGRESS_PORT/graphql" -Method POST -Headers @{"Host"="ldap-manager.localhost"; "Content-Type"="application/json"} -Body $body -UseBasicParsing -TimeoutSec 10 -ErrorAction Stop
        if ($response.StatusCode -eq 200) {
            Write-Host "✅ OK" -ForegroundColor Green
            $result = $response.Content | ConvertFrom-Json
            if ($result.data.login.token) {
                Write-Host "   JWT Token: $($result.data.login.token.Substring(0, 20))..." -ForegroundColor Gray
            }
        } else {
            Write-Host "⚠️  Not ready yet (Status: $($response.StatusCode))" -ForegroundColor Yellow
        }
    } catch {
        Write-Host "⚠️  Not ready yet" -ForegroundColor Yellow
    }

    Write-Host ""
    Write-Host "📌 Access URLs:" -ForegroundColor Cyan
    Write-Host "   Health:  http://$INGRESS_HOST:$INGRESS_PORT/health (Host: ldap-manager.localhost)" -ForegroundColor White
    Write-Host "   GraphQL: http://$INGRESS_HOST:$INGRESS_PORT/graphql (Host: ldap-manager.localhost)" -ForegroundColor White
    Write-Host ""
    Write-Host "💡 Add to C:\Windows\System32\drivers\etc\hosts:" -ForegroundColor Yellow
    Write-Host "   $INGRESS_HOST ldap-manager.localhost" -ForegroundColor White
    Write-Host ""
    Write-Host "🧪 Test with curl or PowerShell:" -ForegroundColor Cyan
    Write-Host "   Invoke-WebRequest -Uri 'http://$INGRESS_HOST:$INGRESS_PORT/graphql' ``" -ForegroundColor Gray
    Write-Host "     -Method POST ``" -ForegroundColor Gray
    Write-Host "     -Headers @{'Host'='ldap-manager.localhost'; 'Content-Type'='application/json'} ``" -ForegroundColor Gray
    Write-Host "     -Body '{\"query\":\"mutation { login(uid: \\\"john.doe\\\", password: \\\"password123\\\") { token } }\"}'" -ForegroundColor Gray
} else {
    Write-Host ""
    Write-Host "⚠️  Ingress Gateway LoadBalancer IP is pending..." -ForegroundColor Yellow
    Write-Host "   Run 'kubectl get svc -n istio-system' to check status" -ForegroundColor White
}

Write-Host ""
Write-Host "✅ Deployment complete!" -ForegroundColor Green
Write-Host ""
Write-Host "📚 Next steps:" -ForegroundColor Cyan
Write-Host "   1. Check logs: kubectl logs -f -l app=ldap-manager -n dev-platform -c ldap-manager" -ForegroundColor White
Write-Host "   2. Check proxy: kubectl logs -f -l app=ldap-manager -n dev-platform -c istio-proxy" -ForegroundColor White
Write-Host "   3. Verify mTLS: istioctl authn tls-check POD_NAME -n dev-platform" -ForegroundColor White
Write-Host "   4. View in Kiali: istioctl dashboard kiali" -ForegroundColor White
Write-Host "   5. Test login credentials: john.doe / password123" -ForegroundColor White
Write-Host ""
