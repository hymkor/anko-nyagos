anko-nyagos
==========

anko-nyagos is the commandline shell based on [nyagos](https://github.com/zetamatta/nyagos) which can be customized by [anko](https://github.com/mattn/anko) instead of Lua.

Sample rc-file: `nyagos.ank` put on the same directory with nyagos.exe.

```
func f(args){
    for s in args{
        println(s)
    }
}

alias("foo",f)
```

The alias named as `foo` is defined by anko.
