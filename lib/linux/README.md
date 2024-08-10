### Guide to build for Linux

#### Install lib

```
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys F6A649124520E5F3
sudo add-apt-repository ppa:ton-foundation/ppa
sudo apt update
sudo apt install ton
```


When you see the successful status of the build, you can find the `libemulator.so` file in the `/opt/homebrew/lib`
folder.

#### Usage

Copy `libemulator.so` to /usr/lib or use environment variable `LD_LIBRARY_PATH`

💡 Full information can be found at github.com/ton-blockchain/packages
