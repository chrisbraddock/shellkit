#!/usr/bin/env python3
"""
copy_claude_plan.py — iTerm2 Python API script

Copies the most recent Claude Code plan block from terminal scrollback
to the macOS system clipboard.

Setup:
  1. Enable Python API: iTerm2 > Settings > General > Magic > Enable Python API
  2. Install Python Runtime if prompted: Scripts > Manage > Install Python Runtime
  3. Deploy to: ~/Library/Application Support/iTerm2/Scripts/AutoLaunch/
  4. Restart iTerm2 (or launch manually from Scripts menu)
  5. Bind hotkey: Settings > Keys > Key Bindings > +
     Action: "Invoke Script Function"
     Function: copy_claude_plan()
"""

import iterm2
import logging
import re
import subprocess

# ── Configuration ─────────────────────────────────────────────────────────────
MAX_LINES = 20000            # Max scrollback lines to scan
SELECT_IN_TERMINAL = True    # Highlight copied block in iTerm2
DEBUG = False                # Enable debug logging to script console

# ── Marker Patterns ──────────────────────────────────────────────────────────
START_MARKER = "Ready to code?"
END_DIVIDER_RE = re.compile(r"\u254c{10,}")   # ╌ (U+254C) repeated 10+ times
END_PHRASE = "Claude has written up a plan"
BORDER_CHARS = frozenset(
    "\u2500\u2501\u2502\u2503\u250c\u250f\u2510\u2513"
    "\u2514\u2517\u2518\u251b\u251c\u2524\u252c\u2534\u253c"
    "\u2550\u2551\u2552\u2553\u2554\u2555\u2556\u2557"
    "\u2558\u2559\u255a\u255b\u255c\u255d\u255e\u255f"
    "\u2560\u2561\u2562\u2563\u2564\u2565\u2566\u2567"
    "\u2568\u2569\u256a\u256b\u256c "
)

# ── Logging ───────────────────────────────────────────────────────────────────
logger = logging.getLogger("copy_claude_plan")
logger.setLevel(logging.DEBUG if DEBUG else logging.WARNING)
if DEBUG:
    _handler = logging.StreamHandler()
    _handler.setFormatter(
        logging.Formatter("%(asctime)s [%(name)s] %(levelname)s: %(message)s"))
    logger.addHandler(_handler)


# ── Search helpers ────────────────────────────────────────────────────────────

def _find_start_marker(lines):
    """Search bottom-up for the last 'Ready to code?' line."""
    for i in range(len(lines) - 1, -1, -1):
        if START_MARKER in lines[i]:
            return i
    return None


def _find_content_start(lines, start_idx):
    """Find where the plan content begins after the header block.

    Starting from the 'Ready to code?' line, scan forward past
    'Here is Claude's plan:', border lines, and blank lines to find
    the first line of actual plan content.
    """
    for i in range(start_idx + 1, min(start_idx + 10, len(lines))):
        stripped = lines[i].strip()
        if not stripped:
            continue
        if "Here is Claude's plan" in stripped:
            continue
        if all(c in BORDER_CHARS for c in stripped):
            continue
        return i
    return start_idx + 1


def _find_end_marker(lines, start_idx):
    """Search forward from start_idx for the end-of-plan marker.

    Returns the index of the end marker line (exclusive boundary),
    or None if not found.
    """
    for i in range(start_idx + 1, len(lines)):
        if END_DIVIDER_RE.search(lines[i]):
            return i
        if END_PHRASE in lines[i]:
            return i
    return None


def _copy_to_clipboard(text):
    """Copy text to macOS clipboard via pbcopy."""
    try:
        subprocess.run(
            ["pbcopy"], input=text.encode("utf-8"), check=True, timeout=5)
        logger.debug("Copied %d bytes to clipboard", len(text))
    except (subprocess.CalledProcessError, subprocess.TimeoutExpired,
            FileNotFoundError) as exc:
        logger.error("pbcopy failed: %s", exc)


# ── Core logic ────────────────────────────────────────────────────────────────

async def _do_copy(connection, session):
    """Fetch scrollback, find plan block, copy to clipboard."""

    # Step 1 — get line geometry in a transaction for consistency
    line_info = await session.async_get_line_info()
    overflow = line_info.overflow
    history = line_info.scrollback_buffer_height
    grid = line_info.mutable_area_height
    total = history + grid

    n = min(total, MAX_LINES)
    first = overflow + (total - n)

    logger.debug(
        "overflow=%d history=%d grid=%d fetching %d from line %d",
        overflow, history, grid, n, first)

    raw = await session.async_get_contents(first, n)
    lines = [line.string for line in raw]

    if not lines:
        logger.info("No scrollback content")
        return

    # Step 2 — find start marker (bottom-up)
    start_idx = _find_start_marker(lines)
    if start_idx is None:
        logger.info("Start marker not found in %d lines", len(lines))
        return

    # Step 3 — find where plan content starts (after header block)
    content_start = _find_content_start(lines, start_idx)

    # Step 4 — find end marker (forward)
    end_idx = _find_end_marker(lines, start_idx)
    if end_idx is None:
        logger.info("No end marker found after start at line %d", start_idx)
        return

    # Step 5 — extract block (end marker line excluded)
    block = lines[content_start:end_idx]
    if not block:
        logger.info("Empty block")
        return

    text = "\n".join(block)
    if not text.endswith("\n"):
        text += "\n"

    logger.debug("Extracted %d lines (%d chars)", len(block), len(text))

    # Step 6 — clipboard
    _copy_to_clipboard(text)

    # Step 7 — optional terminal selection
    if SELECT_IN_TERMINAL:
        await _select_range(session, first, content_start, end_idx)


async def _select_range(session, first_line, block_start, block_end):
    """Highlight the copied region in iTerm2.

    Coordinates are absolute line numbers (first_line + list index).
    CoordRange end point is exclusive — Point(0, end_y) means
    "up to but not including column 0 of that line."
    """
    try:
        abs_start = first_line + block_start
        abs_end = first_line + block_end

        coord = iterm2.CoordRange(
            iterm2.Point(0, abs_start),
            iterm2.Point(0, abs_end))
        windowed = iterm2.WindowedCoordRange(coord)
        sub = iterm2.SubSelection(
            windowed, iterm2.SelectionMode.CHARACTER, False)
        await session.async_set_selection(iterm2.Selection([sub]))

        logger.debug("Selection set y=%d..%d", abs_start, abs_end)
    except Exception as exc:
        logger.warning("Selection failed (non-fatal): %s", exc)


# ── Entry point ───────────────────────────────────────────────────────────────

async def main(connection):
    app = await iterm2.async_get_app(connection)

    @iterm2.RPC
    async def copy_claude_plan(session_id=iterm2.Reference("id")):
        try:
            session = app.get_session_by_id(session_id)
            if not session:
                logger.warning("Session not found: %s", session_id)
                return
            await _do_copy(connection, session)
        except Exception as exc:
            logger.error("copy_claude_plan failed: %s", exc, exc_info=True)

    await copy_claude_plan.async_register(connection)


iterm2.run_forever(main)
