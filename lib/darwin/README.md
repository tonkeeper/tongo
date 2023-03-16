### Guide build for macOS

At the time of writing, there is a bug when building in [the main repository of TON](https://github.com/ton-blockchain/ton/tree/testnet) hasn't been fixed yet, so we're using a fork.

1) `git clone --recurse-submodules -b emulator_vm_verbosity https://github.com/dungeon-master-666/ton.git `
2) `mkdir ton-build && cd ton-build`
3) `cmake ../`
4) `cmake --build . -- target emulator`

When you see the successful status of the build, you can find the `libemulator.dylib` file in the `ton-build/emulator` folder.