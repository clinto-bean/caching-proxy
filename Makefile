.PHONY: dev clean

test:
	@if [ -f ./build/out ]; then \
		rm ./build/out; \
	fi
	go build -o ./build/out
	./build/out serve -p 3001 -s 10 -e 60 -i 6

clean:
	rm -f ./build/out
