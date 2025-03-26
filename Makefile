
hooks:
	cat git-template/hooks/pre-commit >> .git/hooks/pre-commit
	cat git-template/hooks/pre-push >> .git/hooks/pre-push

db:
	docker run --name zchat-db -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust -d postgres
	