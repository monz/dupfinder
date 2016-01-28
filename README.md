# dupfinder
Find duplicate files according to md5 hash

dupfinder takes absolute file paths on stdin, each separated with a newline. It calculates the md5 checksum for each file.
After EOF is read, the summary is printed to stdout. The summary consists of a line with the md5 checksum followed by the
filenames according to the checksum.

# Options
dupfinder takes some options

```
-sumsOnly
        just output the checksums of all files (default 'false')
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

browse all files and subdirectories within `/var/tmp` and print the checksums of all files.
Like the unix tool 'md5sum', but use **10** parallel workers.

```
$ find /var/tmp/ -printf "%p\n" | dupfinder -w 10 -sumsOnly
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

```
$ find /var/tmp/ -printf "%p\n" | dupfinder -sumsOnly
2016/01/27 18:24:04 Worker count 4
8eea72e38a8c03d1932cb505a22c69c7  /var/tmp/file4
5c219e4eef807cb8485e4795fa2ecd1b  /var/tmp/file1
518dea85c42eb48d0db9a5486f9351cc  /var/tmp/file2
38de3a8ad093febed8d7a2e63cdaae37  /var/tmp/file3
```
