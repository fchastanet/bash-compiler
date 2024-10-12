import regex

assertFunctionNamePriority = lambda functionName:str {
  regex.match(functionName, "^([A-Za-z0-9_]+(::)?[A-Za-z0-9_]+)(@[0-9]+)?$")
}

assertFunctionName = lambda functionName:str {
  regex.match(functionName, "^([A-Za-z0-9_]+(::)?[A-Za-z0-9_]+)$")
}
