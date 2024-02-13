dokbuild:
	docker build -t listener .
dokrun:
	docker run -p 8080:1234 listener 
gpull:
	git pull 
gpush:
	git add .
	git commit -m "update"
	git push