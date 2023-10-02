# PDF Import

This documentation describes how a pdf file is imported into the swimresults system.

## General

### Read PDF

All PDF files will be read from the disk and the plain text contents will be parsed as a string.
The processes described in the following part are using this string to work on.

### Extract Events

- split string using event_separator_string ("Wettkampf")
- for all substrings:
  - if string contains one of event_skip_strings (if set), continue
  - if string does not contain one of event_no_skip_strings (if set), continue
  - read following numbers after split until non-number character occurs -> event number
  - search for next number and read until "m"

## Start List

