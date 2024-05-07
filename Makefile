PROJECT_NAME := gociede2k

brun:
	go build .
	./$(PROJECT_NAME)
c2k :
	go build .
	./$(PROJECT_NAME) -mode=c2k
imglib :
	go build .
	./$(PROJECT_NAME) -mode=imglib
vips :
	go build .
	./$(PROJECT_NAME) -mode=vips
