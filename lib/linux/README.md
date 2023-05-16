## Usage

Copy `libemulator.so` to /usr/lib or use environment variable `LD_LIBRARY_PATH`

### Build guide for Linux

    cd lib/linux
    docker build . -t ton-emulator
    docker create --name ton-emulator ton-emulator
    docker cp ton-emulator:/output/libemulator.so .
