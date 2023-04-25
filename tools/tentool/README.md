## Tenhou stats and log downloader
Download stat files and log files to work with them on your local computer.

Download tool by typing:
```
go get github.com/dnovikoff/tenhou/tools/tentool
```

Initialize the stats repo in current working dir.
```
tentool stats init
```

I suggest that you first download some archives from my [Yandex.Disk](https://yadi.sk/d/uOv87aVsd-l-3A).
The files will be downloaded to `./tenhou/stats` folder.
That could you reduce you downloading time. 
Do this by typing:
```
tentool stats yadisk
```

Download stat files from tenhou.net.
The files will be downloaded to `./tenhou/stats` folder.
Repeat this action when you need to get updates.
```
tentool stats download
```

Initialize the logs repo in current working dir.
```
tentool logs init
```

I suggest that you first download some prebuild zip files with logs from my [Yandex.Disk](https://yadi.sk/d/FIIkaucSNjR3Kw).
That would be sure times faster, than downloading all logs from tenhou one by one.
The files will be downloaded to `./tenhou/logs` folder.
```
tentool logs yadisk
```

Collect all log ids from stat files by typing.
Repeat this action, after next call of `tentool stats download`.
```
tentool logs update
```

Download log files from collected log ids by typing
```
tentool logs download
```

Alternatevly you can init with makefile
```
make init
```

And update with
```
make download
```

Now you have full database of phoenix logs on your machine.
Consider reading `tentool stats --help` and `tentool logs --help` on more commands and flags.