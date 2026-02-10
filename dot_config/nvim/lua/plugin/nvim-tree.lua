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

-- Auto-close nvim-tree when it's the last window remaining
vim.api.nvim_create_autocmd("QuitPre", {
    callback = function()
        local tree_wins = {}
        local floating_wins = {}
        local wins = vim.api.nvim_list_wins()
        for _, w in ipairs(wins) do
            local bufname = vim.api.nvim_buf_get_name(vim.api.nvim_win_get_buf(w))
            if bufname:match("NvimTree_") ~= nil then
                table.insert(tree_wins, w)
            end
            if vim.api.nvim_win_get_config(w).relative ~= "" then
                table.insert(floating_wins, w)
            end
        end
        if #wins - #floating_wins - #tree_wins == 1 then
            for _, w in ipairs(tree_wins) do
                vim.api.nvim_win_close(w, true)
            end
        end
    end,
})

-- Keymaps
vim.keymap.set("n", "<C-n>", "<cmd>NvimTreeToggle<CR>", { desc = "Toggle file explorer" })
vim.keymap.set("n", "<leader>e", "<cmd>NvimTreeFocus<CR>", { desc = "Focus file explorer" })
vim.keymap.set("n", "<leader>f", "<cmd>NvimTreeFindFile<CR>", { desc = "Find current file in tree" })
