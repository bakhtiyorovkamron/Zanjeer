dokbuild:
	docker build -t listener .
dokrun:
	docker run -p 8080:1234 listener .