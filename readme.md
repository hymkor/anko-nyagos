anko-nyagos
==========

anko-nyagos is the commandline shell based on [nyagos](https://github.com/zetamatta/nyagos) which can be customized by [anko](https://github.com/mattn/anko) instead of Lua.

Sample rc-file: nyagos.ank which is put on the same directory with nyagos.exe.

```
func f(value){
    println(value)
}

alias("foo",f)
```
The alias `foo` is defined by anko.



