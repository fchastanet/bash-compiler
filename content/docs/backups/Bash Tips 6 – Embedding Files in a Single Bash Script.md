---
title: Bash Tips 6 – Embedding Files in a Single Bash Script
description: Techniques for embedding files into a single bash script using base64 encoding and process substitution
weight: 20
categories: [documentation]
tags: [bash, scripting, embedding files, base64 encoding, process substitution]
backup:
  author: Michał Zieliński
  authorUrl: https://blog.tratif.com/author/mzielinski/
  originalUrl: https://blog.tratif.com/2023/02/17/bash-tips-6-embedding-files-in-a-single-bash-script/
creationDate: 2023-02-17
lastUpdated: 2026-02-27
version: '1.0'
---

Scripts that utilize multiple files are not easy to distribute. We usually distribute those as archives and rely on the
end user to unpack and run them from a predetermined location. To improve the experience we can instead prepare a single
script with other files embedded inside it.

Here are the goals:

1. The script should consist of a single file, making it easy to distribute.
2. The script should be copy-paste-able between systems and different editors, even if multiple hops are required.
3. Files being embedded can be binary files i.e. can contain non-printable characters.
4. The first requirement implies that we should somehow store the contents of other files in our main script. The second
   requires us to avoid non-printable characters, as they tend to cause problems when performing a copy-paste operation.
   Especially when we are talking about sending such characters over messaging programs.

# Encoding

The solution to the second and third problems is a binary-to-text encoding which encodes an array of bytes into a text
constant of printable characters. And the most commonly used encoding scheme is base64. Utils to encode to and from
base64 are included in most Linux distributions out-of-the-box.

Let’s transform a file, logging setup script, into base64 encoded text:

```bash
base64 -w 0 includes/logging.sh
```

prints

```text
TE9HRklMRT0iJHsxOi1zY3JpcHQubG9nfSIKZXhlYyAzPiYxIDE+IiRMT0dGSUxFIiAyPiYxCnRyYXAgImVjaG8gJ0VSUk9SOiBBbiBlcnJvciBvY2N1cnJlZCBkdXJpbmcgZXhlY3V0aW9uLCBjaGVjayBsb2cgJExPR0ZJTEUgZm9yIGRldGFpbHMuJyA+JjMiIEVSUgp0cmFwICd7IHNldCAreDsgfSAyPi9kZXYvbnVsbDsgZWNobyAtbiAiWyQoZGF0ZSAtSXMpXSAgIjsgc2V0IC14JyBERUJVRwo=
```

which, when decoded

```bash
base64 -d <<<TE9HRklMRT0iJHsxOi1zY3JpcHQubG9nfSIKZXhlYyAzPiYxIDE+IiRMT0dGSUxFIiAyPiYxCnRyYXAgImVjaG8gJ0VSUk9SOiBBbiBlcnJvciBvY2N1cnJlZCBkdXJpbmcgZXhlY3V0aW9uLCBjaGVjayBsb2cgJExPR0ZJTEUgZm9yIGRldGFpbHMuJyA+JjMiIEVSUgp0cmFwICd7IHNldCAreDsgfSAyPi9kZXYvbnVsbDsgZWNobyAtbiAiWyQoZGF0ZSAtSXMpXSAgIjsgc2V0IC14JyBERUJVRwo=
```

prints the contents of the file:

```bash
LOG_FILE="${1:-script.log}"
exec 3>&1 1>"$LOG_FILE" 2>&1
trap "echo 'ERROR: An error occurred during execution, check log $LOG_FILE for details.' >&3" ERR
trap '{ set +x; } 2>/dev/null; echo -n "[$(date -Is)]  "; set -x' DEBUG
```

this can be redirected to a file, that can later be used:

```bash
base64 -d >/tmp/logging.sh <<<TE9HRklMRT0iJHsxOi1zY3JpcHQubG9nfSIKZXhlYyAzPiYxIDE+IiRMT0dGSUxFIiAyPiYxCnRyYXAgImVjaG8gJ0VSUk9SOiBBbiBlcnJvciBvY2N1cnJlZCBkdXJpbmcgZXhlY3V0aW9uLCBjaGVjayBsb2cgJExPR0ZJTEUgZm9yIGRldGFpbHMuJyA+JjMiIEVSUgp0cmFwICd7IHNldCAreDsgfSAyPi9kZXYvbnVsbDsgZWNobyAtbiAiWyQoZGF0ZSAtSXMpXSAgIjsgc2V0IC14JyBERUJVRwo=
```

Using base64 allows us to store binary files as an easy to work with text. I have used exactly this mechanism to prepare
a script that would store and activate binary licenses for an external proprietary system that we have used in one of
our projects.

It is worth noting that in the case of shell scripts, base64 encoding provides a safety layer preventing us from
accidental execution. If we were to use a here-document to try to achieve the same functionality would have to account
for variable expansion:

```bash
cat >/tmp/logging.sh <<EOC
LOG_FILE="${1:-script.log}"
exec 3>&1 1>"$LOG_FILE" 2>&1
trap "echo 'ERROR: An error occurred during execution, check log $LOG_FILE for details.' >&3" ERR
trap '{ set +x; } 2>/dev/null; echo -n "[$(date -Is)]  "; set -x' DEBUG
EOC
```

This code does not work, the resulting file has variables expanded, as all variables inside here-documents are expanded:

(contents of /tmp/logging.sh file created by running the command above)

```bash
LOG_FILE="script.log"
exec 3>&1 1>"" 2>&1
trap "echo 'ERROR: An error occurred during execution, check log  for details.' >&3" ERR
trap '{ set +x; } 2>/dev/null; echo -n "[2023-01-24T13:36:23+01:00]  "; set -x' DEBUG
```

# Utilizing process substitution

We do not have to create a temporary file that we later have to clean up. If a file is only to be read once, we can
utilize process substitution. It allows referencing an output (or input) of a process as a file that can be accessed.
Let’s see that in an example:

```bash
base64 -d <<<SGVsbG8gV29ybGQhCg==
```

this command prints ‘Hello World!’ on standard output, nothing special.

```bash
<(base64 -d <<<SGVsbG8gV29ybGQhCg==)
```

this returns a path to a file, that, when read, returns the output of the command inside the `<(...)`.

```bash
echo <(base64 -d <<<SGVsbG8gV29ybGQhCg==)
```

prints

```bash
/proc/self/fd/11
```

a file path. This is a file created thanks to using process substitution. If we read the file we get the output of our
base64 which contains decoded message. Let’s replace the echo in the call with cat to see the file contents:

```bash
cat <(base64 -d <<<SGVsbG8gV29ybGQhCg==)
```

`<(base64 -d <<<SGVsbG8gV29ybGQhCg==)` gets transformed into a file path, and cat performs read on this file and prints
its contents on the console:

```text
Hello World!
```

Quite a sophisticated way to print a simple message. With the mechanism explained, let’s proceed with utilizing it in a
script:

```bash
source <(base64 -d <<<TE9HRklMRT0iJHsxOi1zY3JpcHQubG9nfSIKZXhlYyAzPiYxIDE+IiRMT0dGSUxFIiAyPiYxCnRyYXAgImVjaG8gJ0VSUk9SOiBBbiBlcnJvciBvY2N1cnJlZCBkdXJpbmcgZXhlY3V0aW9uLCBjaGVjayBsb2cgJExPR0ZJTEUgZm9yIGRldGFpbHMuJyA+JjMiIEVSUgp0cmFwICd7IHNldCAreDsgfSAyPi9kZXYvbnVsbDsgZWNobyAtbiAiWyQoZGF0ZSAtSXMpXSAgIjsgc2V0IC14JyBERUJVRwo=)
```

This is a base64 encoded logging script from the previous section, which we can easily source without creating a
temporary file.

# Embedding entire directory tree

With the tools we have at our disposal presented in previous sections, let’s try to create a single portable script from
a complex multi-file setup. We will use the code presented in a previous article of mine on templating.

First, let’s change our working directory to the one that contains the main script script-3.1.sh and create a compressed
archive of the contents of the entire project using tar and gunzip, which we base64 encode:

```bash
cd ... # path to the project

# .
# ├── includes
# │   ├── gatheringFacts.sh
# │   ├── logging.sh
# │   └── templating.sh
# ├── script-3.1.sh
# ├── templates
# │   └── config.yml
# ├── utils
# │   └── getIp.sh

tar -cz -O . | base64 -w 0
```

Here is our encoded archive:

```base64
H4sIAAAAAAAAA+1Z+2/iRhDmZ0v5H+Zc98C64iePiJRc6R25IuWBOKhUJShy7AVWMba1a0LSHP97Z21eIdfSVoGoV39IRp79dmbx7IxnFk3P7RwGoloui2+zWjbWvxfImSW7YlfLlmGZOcPEayUH5d0vLZeb8NhhALkbEofBX/C2jf9HoemTmPp8p7vg7/u/VClbNvrftKtG5v99YOH/IYlbkcZHu7Cxzf8V217EfwV9j/63qlUrB8YuFrOJ/7n/v3uj39BAv3H4SJpwZ0gKKjxKgCDuKAS5zWgQcygMKOOxCq02OJ7HCOfgcB661ImJB1MajyAg8TRkt4B8wgaOS+COOjAdUXeUypAAMXMGA+pCHAKN4FBLPkA5sHCCmjRZmknSUkNdKSArGQLcoEv+FxgyEkGxHUK+8P7HukfuQL2aviu8r4OaV6X5EtPp4gb4KJyCoCmr5a1pKY7BTFVRsUr1ynt3pT25oGpdqJboAC4vofg7yMrcjAz9/hHEIxKsPbdmp3PRqcGJQ318QPhzxfqRg09QhuO3Vkq9pzGY0oBKUjKrGKxplfbif03nLqNRXLQ1c0fhvzX+bXMz/1t2pZTF/z6wHv9Su9H9pa6Ia00pRFNPTV8OEg8nDAOGBq4/8QjX/XA4pMEQ9wvI6faxNZTJz4hDB7c8ZpDhiePGHPnPGDEZR74Tp9okaX5LQD5pfOh+lmEh4LobBgM61B7GPujxONLHDyuJtJ9g+Qah6asnvCsb/7z+sywrq//3gnX/r+LpZW0kNV6p9Kf+t4zqhv9LllnO8v8+wAm7I6x2kNQj89KjBsrnZufXZue61U4HopDFNTg0Do0D6UCiHgliOqA4b8X82Dzvtk5azc6B5EYT6q2GPrR712eN857I6L1Owj3AWjOZ/dhD0uwgS9+vBU1fvop3ZuNf5H/DzvL/XrDm/1VV98I2tuV/u2yu+v9KVeR/A7dBlv/3gNOLTyet02ZdVh7NWhG3gCjkZ7JE7okL9vFbE8xjWZmzZLBQImETH4GcdKz5eZ/bCIAwFjIIXXfCGLa83kTU/SD0YBMRBj+AOyLuLaB6UAqMOH6EzQEsVKswwNkeibFh5loeW2RbBlSeGss/Asf++d39EcxwDTp28now8f0jWPbNl0rBE41DscXVPoB8lMwo3ufhY/Pn3qfsDfN1rMX/s2btpWxsrf8qm/Vf2baNLP73AY+4vsMwbBqQdNxScr3MLwvAfL+uFBanw+rm8LLqS2ijkMeBMyabtK9UgKna5PxNnL3dkcAL2TX1QI9Y6OqihAwGIXwBzB5Q9KAGxQFYeH/vsCFXX/uxfTNYi/8nRzEvaWNb/Jdsa7P/K9tZ/7cXLLp/cewPSa/X6HQav12fN85EUWDKiazbPGufNrrN6+SEUFasVHzR67Z73YXQlqVE6oeu44t3cqIJB1Ya02niTX9LHoAG2AC+SUYvf+rPjsALk/H0aFz0nKAgr648phxx058lFC8MSGrNdWKQlScLlDFNkOCOT264GCtE4j+MAeSV7znkUbBmU1YFOSYExWu/RvwN8dqeyZAhQ4YMGTJk2A3+APiee2EAKAAA
```

We now can create a wrapper script, that will create a temporary directory, change the working directory to it and
extract the files there. Then it would run the main script.

```bash
#!/bin/bash
cd $(mktemp -d)
tar -xzf <(base64 -d <<<H4sIAAAAAAAAA+1Z+2/iRhDmZ0v5H+Zc98C64iePiJRc6R25IuWBOKhUJShy7AVWMba1a0LSHP97Z21eIdfSVoGoV39IRp79dmbx7IxnFk3P7RwGoloui2+zWjbWvxfImSW7YlfLlmGZOcPEayUH5d0vLZeb8NhhALkbEofBX/C2jf9HoemTmPp8p7vg7/u/VClbNvrftKtG5v99YOH/IYlbkcZHu7Cxzf8V217EfwV9j/63qlUrB8YuFrOJ/7n/v3uj39BAv3H4SJpwZ0gKKjxKgCDuKAS5zWgQcygMKOOxCq02OJ7HCOfgcB661ImJB1MajyAg8TRkt4B8wgaOS+COOjAdUXeUypAAMXMGA+pCHAKN4FBLPkA5sHCCmjRZmknSUkNdKSArGQLcoEv+FxgyEkGxHUK+8P7HukfuQL2aviu8r4OaV6X5EtPp4gb4KJyCoCmr5a1pKY7BTFVRsUr1ynt3pT25oGpdqJboAC4vofg7yMrcjAz9/hHEIxKsPbdmp3PRqcGJQ318QPhzxfqRg09QhuO3Vkq9pzGY0oBKUjKrGKxplfbif03nLqNRXLQ1c0fhvzX+bXMz/1t2pZTF/z6wHv9Su9H9pa6Ia00pRFNPTV8OEg8nDAOGBq4/8QjX/XA4pMEQ9wvI6faxNZTJz4hDB7c8ZpDhiePGHPnPGDEZR74Tp9okaX5LQD5pfOh+lmEh4LobBgM61B7GPujxONLHDyuJtJ9g+Qah6asnvCsb/7z+sywrq//3gnX/r+LpZW0kNV6p9Kf+t4zqhv9LllnO8v8+wAm7I6x2kNQj89KjBsrnZufXZue61U4HopDFNTg0Do0D6UCiHgliOqA4b8X82Dzvtk5azc6B5EYT6q2GPrR712eN857I6L1Owj3AWjOZ/dhD0uwgS9+vBU1fvop3ZuNf5H/DzvL/XrDm/1VV98I2tuV/u2yu+v9KVeR/A7dBlv/3gNOLTyet02ZdVh7NWhG3gCjkZ7JE7okL9vFbE8xjWZmzZLBQImETH4GcdKz5eZ/bCIAwFjIIXXfCGLa83kTU/SD0YBMRBj+AOyLuLaB6UAqMOH6EzQEsVKswwNkeibFh5loeW2RbBlSeGss/Asf++d39EcxwDTp28now8f0jWPbNl0rBE41DscXVPoB8lMwo3ufhY/Pn3qfsDfN1rMX/s2btpWxsrf8qm/Vf2baNLP73AY+4vsMwbBqQdNxScr3MLwvAfL+uFBanw+rm8LLqS2ijkMeBMyabtK9UgKna5PxNnL3dkcAL2TX1QI9Y6OqihAwGIXwBzB5Q9KAGxQFYeH/vsCFXX/uxfTNYi/8nRzEvaWNb/Jdsa7P/K9tZ/7cXLLp/cewPSa/X6HQav12fN85EUWDKiazbPGufNrrN6+SEUFasVHzR67Z73YXQlqVE6oeu44t3cqIJB1Ya02niTX9LHoAG2AC+SUYvf+rPjsALk/H0aFz0nKAgr648phxx058lFC8MSGrNdWKQlScLlDFNkOCOT264GCtE4j+MAeSV7znkUbBmU1YFOSYExWu/RvwN8dqeyZAhQ4YMGTJk2A3+APiee2EAKAAA)
./script-3.1.sh
```

# Summary

Using the techniques I have shown in this article we achieved the goals stated at the beginning of the article:

The script should consist of a single file, making it easy to distribute The script should be copy-paste-able between
systems and different editors, even if connecting via multiple jumps Files being embedded can be binary files i.e. can
contain non-printable characters A complex multi-file directory structure has been transformed into a copy-paste-able,
few lines long, script. Of course, this practice should be used when necessary and should be avoided if possible. The
contents, hidden behind encoding and compression, are completely obfuscated. The user has no idea what such script does
and has to reverse engineer the process to find out. I have personally used this when working with certain client
servers that were reachable via multiple jump hosts and had no internet access, where copy-pasting a single script was a
big time saver over transferring the files.
