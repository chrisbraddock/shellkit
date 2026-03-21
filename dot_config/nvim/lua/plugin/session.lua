-- Auto-save nvim session on exit (for tmux-resurrect restore)
-- Sessions are saved per working directory to ~/.local/state/nvim/sessions/

local session_dir = vim.fn.stdpath("state") .. "/sessions"

local function get_session_file()
    local cwd = vim.fn.getcwd()
    local name = cwd:gsub("/", "%%")
    return session_dir .. "/" .. name .. ".vim"
end

vim.api.nvim_create_autocmd("VimLeavePre", {
    group = vim.api.nvim_create_augroup("AutoSession", { clear = true }),
    callback = function()
        -- Only save if we have real buffers open (not just empty scratch)
        local bufs = vim.tbl_filter(function(b)
            return vim.bo[b].buflisted and vim.bo[b].buftype == ""
        end, vim.api.nvim_list_bufs())
        if #bufs == 0 then return end

        vim.fn.mkdir(session_dir, "p")
        vim.cmd("mksession! " .. vim.fn.fnameescape(get_session_file()))
    end,
})
