
The idea of this tool came to me as I was re-reading the original handwriting of Matthew 25:32 (recently recovered and restored):

![image](logshepherd.png "Matthew 25:32")

---
### __And all the log entries shall be gathered before him, and he shall separate them one from another, as a shepherd separates the sheep from the goats.__
---
<br><br>

# Logshepherd

The idea of this tool is - we replace the regular log files with named pipes (aka fifo) and we put the Logshepherd on the receiving end of these pipes.

Now we can do this - we can set the maximum verbosity in all of log producers and let the Logshepherd process all of the messages coming in on the pipes. Based on a very simple filter (plain text match, no fancy stuff like regex) the Logshepherd will output only the messages of the desired log level (e.g., errors) into the corresponding regular log files.

So far, maybe seems useless, but there is more. The combination of "show" and "hide" filters actually allows to catch very specific messages in the logs. For example, we want the usual verbosity + any messages in any verbosity which occur around midnight - can be done.

Logshepherd can also optionally output the whole incoming stream (with a timestamp) into a separate file. It will also rotate both the regular log and the full output when they reach predefined size.

Logshepherd has a configurable buffer for each log. This buffer can be flushed into the regular log, if the message matches one of the trigger filters ("ERROR", for example). This allows us to look deeper into error conditions without consuming gigabytes of disk space for full verbosity logging.

## How to Use

To run and read the configuration from `logshepherd.yaml` in the same directory:

    $ ./logshepherd

The path to a different configuration file can be given as an optional parameter:

    $ ./logshepherd /etc/logshepherd.yaml

The configuration in `logshepherd.yaml` is pretty self-explanatory.

## TODO

- configuration reload on signal
- __dump_buffer__ AFTER the trigger
- also __dump_timeout__ for the AFTER part
- rotated logs compression
