### Guide build for macOS


1) `git clone --recurse-submodules -b emulator-optional-config https://github.com/dungeon-master-666/ton.git `
2) `mkdir ton-build && cd ton-build`
3) `cmake ../`
4) `cmake --build . -- target emulator`

When you see the successful status of the build, you can find the `libemulator.dylib` file in the `ton-build/emulator` folder.