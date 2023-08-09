### Guide build for macOS


1) `git clone --recurse-submodules -b master https://github.com/ton-blockchain/ton.git`
2) `mkdir build && cd build`
3) `cmake ..`
4) `cmake --build . -- target emulator`

When you see the successful status of the build, you can find the `libemulator.dylib` file in the `build/emulator` folder.