package = "gleamold-lua"
version = "scm-02"

source = {
    url = "git://github.com/chrislusf/gleamold.git",
}

description = {
    summary = "Lua libraries to work with Gleamold",
    homepage = "https://github.com/chrislusf/gleamold",
    license = "MIT/X11",
    maintainer = "Chris Lu <chris.lu@gmail.com>",
    detailed = [[
Gleamold-Lua a high-performance library for Gleamold.
It works with Luajit and Lua.
]]
}

dependencies = {
    "lua"
}

build = {
    type = "builtin",
    modules = {
        MessagePack = "script/MessagePack.lua",
    },
}