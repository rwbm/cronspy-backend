# run tests on all folders except 'vendor'
go test -race -coverprofile=profile.out -covermode=atomic $(go list ./...)

if [ $? -eq 0 ]; then
    go tool cover -html=profile.out
fi