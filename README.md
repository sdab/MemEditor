# MemEditor
A simple linux process memory editor in go.

## How to use
Run the process you'd like to scan/edit. Figure out its pid using ps or the like. Run `sudo MemEditor --pid $PID` to start scanning.

The editor runs a read-eval-loop interactive prompt. The interactive prompt includes the number of tracked addresses. (with `%d>:` where %d is the number of tracked addresses).

The commands you can use in the prompt are:
*  `scan val` - scans the current values of all tracked addresses and
   filters the tracked addresses by value. Scans the whole of mapped
   memory if there are no tracked addresses (such as on startup or
   after a reset).
*  `list` - lists all tracked addresses and their last values.
*  `update` - scans the current values of all tracked addresses.
*  `set addr val` - Writes val to address addr.
*  `setall val` - Writes val to all tracked addresses.
*  `reset` - Removes all tracked addresses. The next scan will read
   all of mapped memory.
*  `help` - prints the commands
