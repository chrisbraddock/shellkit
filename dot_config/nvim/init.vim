" ============================================================
" Neovim Configuration
" ============================================================
" Sources shared vim config from ~/.vimrc
" ============================================================

" Source shared vim configuration
source ~/.vimrc

" ------------------------------------------------------------
" Neovim-specific settings
" ------------------------------------------------------------

" Use true colors
set termguicolors

" Live substitution preview
set inccommand=split

" Load Lua configuration (plugins, nvim-tree, etc.)
lua require('init')
