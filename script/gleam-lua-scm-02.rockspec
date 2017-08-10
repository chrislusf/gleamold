package = "gleamold-lua"
version = "scm-02"

source = {
    url = "git://github.com/chrislusf/gleamold.git",
}

description = {
    summary = "Lua libraries to work with Gleam",
    homepage = "https://github.com/chrislusf/gleamold",
    license = "MIT/X11",
    maintainer = "Chris Lu <chris.lu@gmail.com>",
    detailed = [[
Gleam-Lua a high-performance library for Gleam.
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