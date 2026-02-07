-- Plugin manager (lazy.nvim) - self-bootstrapping

local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"
if not vim.loop.fs_stat(lazypath) then
    vim.fn.system({
        "git",
        "clone",
        "--filter=blob:none",
        "https://github.com/folke/lazy.nvim.git",
        "--branch=stable",
        lazypath,
    })
end
vim.opt.rtp:prepend(lazypath)

require("lazy").setup({
    {
        "nvim-tree/nvim-tree.lua",
        dependencies = { "nvim-tree/nvim-web-devicons" },
        config = function()
            require("plugin.nvim-tree")
        end,
    },
}, {
    install = {
        colorscheme = { "default" },
    },
    checker = {
        enabled = false,
    },
    change_detection = {
        notify = false,
    },
})
