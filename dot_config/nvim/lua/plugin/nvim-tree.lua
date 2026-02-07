-- nvim-tree: file explorer sidebar

-- Disable netrw (avoids conflicts with nvim-tree)
vim.g.loaded_netrw = 1
vim.g.loaded_netrwPlugin = 1

require("nvim-tree").setup({
    sort = {
        sorter = "case_sensitive",
    },
    view = {
        width = 30,
        side = "left",
    },
    renderer = {
        group_empty = true,
        icons = {
            show = {
                file = true,
                folder = true,
                folder_arrow = true,
                git = true,
            },
        },
    },
    filters = {
        dotfiles = false,
        custom = { ".git", "node_modules", ".DS_Store" },
    },
    git = {
        enable = true,
        ignore = false,
    },
    actions = {
        open_file = {
            quit_on_open = false,
            resize_window = true,
        },
    },
    update_focused_file = {
        enable = true,
    },
})

-- Auto-open for directories or no-arg launches; skip for single files
vim.api.nvim_create_autocmd("VimEnter", {
    callback = function(data)
        local is_directory = vim.fn.isdirectory(data.file) == 1
        local no_args = data.file == "" and vim.bo[data.buf].buftype == ""

        if is_directory then
            vim.cmd.cd(data.file)
            require("nvim-tree.api").tree.open()
        elseif no_args then
            require("nvim-tree.api").tree.open()
        end
    end,
})

-- Keymaps
vim.keymap.set("n", "<C-n>", "<cmd>NvimTreeToggle<CR>", { desc = "Toggle file explorer" })
vim.keymap.set("n", "<leader>e", "<cmd>NvimTreeFocus<CR>", { desc = "Focus file explorer" })
vim.keymap.set("n", "<leader>f", "<cmd>NvimTreeFindFile<CR>", { desc = "Find current file in tree" })
