from wasmtime import Store, Module, Instance

store = Store()
module = Module.from_file(store.engine, 'gcd.wat')
instance = Instance(store, module, [])
gcd = instance.exports(store)['gcd']
print("gcd(27, 6) = %d" % gcd(store, 27, 6))

