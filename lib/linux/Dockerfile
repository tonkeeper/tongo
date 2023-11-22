FROM docker.io/library/ubuntu:20.04 AS emulator-builder
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get -y install build-essential git cmake clang libgflags-dev zlib1g-dev libssl-dev libreadline-dev libmicrohttpd-dev pkg-config libgsl-dev python3 python3-dev python3-pip libsecp256k1-dev libsodium-dev
RUN git clone --recurse-submodules -b testnet https://github.com/ton-blockchain/ton.git
RUN mkdir build && (cd build && cmake ../ton -DCMAKE_BUILD_TYPE=Release && cmake --build . --target emulator -j 6)
RUN mkdir /output && cp build/emulator/libemulator.so /output
