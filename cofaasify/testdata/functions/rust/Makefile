# MIT License
#
# Copyright (c) 2021 Michal Baczun and EASE lab
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

all: all-image protos

all-image: producer-image consumer-image

all-image-push: producer-image-push consumer-image-push driver-image-push

helloworld_protos := producer/helloworld.proto
prodcon_protos := consumer/prodcon.proto producer/prodcon.proto
protos: $(helloworld_protos) $(prodcon_protos)

clean:
	rm -f ./producer-bin ./consumer-bin


producer-image: Dockerfile producer/src/main.rs $(helloworld_protos) $(prodcon_protos)
	DOCKER_BUILDKIT=1 docker build \
	--tag cofaas/producer-rust:latest \
	--build-arg target_arg=producer \
	-f Dockerfile .

producer-image-push: producer-image
	docker push vhiveease/chained-functions-serving-producer:latest

consumer-image: Dockerfile consumer/src/main.rs $(prodcon_protos)
	DOCKER_BUILDKIT=1 docker build \
	--tag cofaas/consumer-rust:latest \
	--build-arg target_arg=consumer \
	-f Dockerfile .

consumer-image-push: consumer-image
	docker push vhiveease/chained-functions-serving-consumer:latest

up:
	podman-compose up

.PHONY: up
