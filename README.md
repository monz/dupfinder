# dupfinder
Find duplicate files according to md5 hash

dupfinder takes absolute file paths on stdin, each separated with a newline. It calculates the md5 checksum for each file.
After EOF is read, the summary is printed to stdout. The summary consists of a line with the md5 checksum followed by the
filenames according to the checksum.

# Options
dupfinder takes one option, the worker count

```
-w int
    	count of parallel md5sum workers (default 4)
```

# Examples

## Linux
browse all files and subdirectories within `/var/tmp` and search for duplicate files. Save duplicates in a text file.
Use **10** md5sum workers

```
$ find /var/tmp/ -printf "%p\n" | dupfinder -w 10 > duplicates.txt
```

## Windows
browse all files and subdirectories within `C:\tmp` and search for duplicate files. Save duplicates in a text file

```
> dir C:\tmp /S /B" | dupfinder.exe > duplicates.txt
```

## Output
```
$ find /var/tmp/ -printf "%p\n" | dupfinder
2015/11/22 20:07:38 Worker count 4
Checksum d41d8cd98f00b204e9800998ecf8427e:
  /var/tmp/file4
  /var/tmp/file1
  /var/tmp/file2
  /var/tmp/file3
  /var/tmp/file5
  /var/tmp/file6
```
