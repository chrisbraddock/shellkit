-- Treesitter: syntax highlighting for common languages

local langs = {
    "bash",
    "css",
    "dockerfile",
    "html",
    "javascript",
    "json",
    "lua",
    "markdown",
    "markdown_inline",
    "python",
    "toml",
    "typescript",
    "vim",
    "vimdoc",
    "yaml",
}

-- Install parsers
require("nvim-treesitter").setup({
    ensure_installed = langs,
})

-- Enable treesitter highlighting for all buffers with a parser available
vim.api.nvim_create_autocmd("FileType", {
    callback = function(args)
        pcall(vim.treesitter.start, args.buf)
    end,
})
