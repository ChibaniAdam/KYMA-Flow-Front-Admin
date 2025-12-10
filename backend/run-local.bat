@echo off
echo Loading environment variables...
set LDAP_URL=ldap://localhost:30000
set LDAP_BASE_DN=dc=devplatform,dc=local
set LDAP_BIND_DN=cn=admin,dc=devplatform,dc=local
set LDAP_BIND_PASSWORD=admin123
set JWT_SECRET=dev-secret-key-change-in-production-12345678
set PORT=8080
set METRICS_PORT=9090
set ENVIRONMENT=development
set LOG_LEVEL=debug
set LDAP_POOL_SIZE=5
set STARTING_UID=10000
set STARTING_GID=10000

echo Starting LDAP Manager Service...
echo.
echo Service will be available at:
echo   - GraphQL API: http://localhost:8080/graphql
echo   - Health Check: http://localhost:8080/health
echo   - Readiness: http://localhost:8080/ready
echo   - Metrics: http://localhost:9090/metrics
echo.

go run cmd/server/main.go
