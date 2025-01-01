# sponge-go

Pure Go implementation of sponge(1) in moreutils.

> [!WARNING]
> The buffer size and signal handling are roughly implemented. However, In most cases you don't need to worry.


# Usage

```
sponge [-a] <file>: soak up all input from stdin and write it to <file>
```

```
sed '...' file | grep '...' | sponge [-a] file
```



# References

* [moreutils](https://joeyh.name/code/moreutils/)
* [sponge(1) — moreutils — Debian testing — Debian Manpages](https://manpages.debian.org/testing/moreutils/sponge.1.en.html)
* [moreutils/sponge.c at master · stigtsp/moreutils · GitHub](https://github.com/stigtsp/moreutils/blob/master/sponge.c)
* [Redirections (Bash Reference Manual)](https://www.gnu.org/software/bash/manual/html_node/Redirections.html)


