import run_lambda
import imported_source

# sourcefile.py imports imported_source.py and changes a variable in it
# we want to make sure that change is only visible INSIDE of sourcefile.py, and not out here

print 'In \'parent process\', imported_source.c has value: ' + str(imported_source.c)

print 'Running lambda with protection'
run_lambda.executeFileWithProtection('sourcefile.py')

print 'In \'parent process\', imported_source.c has value: ' + str(imported_source.c)

print 'Running lambda in the unprotected way'
run_lambda.executeFileWithoutProtection('sourcefile.py')

print 'In \'parent process\', imported_source.c has value: ' + str(imported_source.c)
