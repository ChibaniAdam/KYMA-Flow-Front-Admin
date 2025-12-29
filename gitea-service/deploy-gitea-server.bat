@echo off
echo Deploying Gitea Server...
echo.

echo Step 1: Creating dev-platform namespace if not exists...
kubectl create namespace dev-platform 2>nul
echo.

echo Step 2: Deploying Gitea Server...
kubectl apply -f .\k8s\00-gitea-server.yaml
echo.

echo Step 3: Waiting for Gitea to be ready...
echo This may take a few minutes on first deployment...
kubectl wait --for=condition=ready pod -l app=gitea -n dev-platform --timeout=300s
if errorlevel 1 (
    echo Gitea not ready yet, check status manually
    echo Run: kubectl get pods -n dev-platform -l app=gitea
) else (
    echo Gitea is ready!
)
echo.

echo Step 4: Verifying deployment...
echo.
echo === Gitea Pod ===
kubectl get pods -n dev-platform -l app=gitea
echo.
echo === Gitea Service ===
kubectl get svc gitea -n dev-platform
echo.

echo ========================================
echo Gitea Server deployment complete!
echo ========================================
echo.
echo Access Gitea at: http://localhost:30009
echo SSH access on: localhost:30010
echo.
echo First-time setup:
echo 1. Open http://localhost:30009 in your browser
echo 2. Complete the initial setup wizard
echo 3. Create an admin user
echo 4. Generate an access token for the gitea-service
echo.
echo To check logs:
echo   kubectl logs -f -l app=gitea -n dev-platform
echo.
