# MP4toWEBmCLI
Simmilar to te other repo but its not a service but a cli tool. In fatc a very simple one.
Code is ok but lacks features and for some reason I didn't use the ffmpeg librarly for go.

## Idea

MP4 files didn't work in Discord when downloaded from the web.



## Installation and usage

Just build the tool with go compiler and run the exe.
The tool will show mp4 files in Download folder that are less than 1 day old.
The convertion will be saved in the Downloads folder.


```sh
go build main.go
run the exe
```

Select the file and wait
If file is very big it will convert a long time.
