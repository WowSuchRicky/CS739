import run_lambda
import imported_source
import numpy

# sourcefile.py imports imported_source.py and changes a variable in it
# we want to make sure that change is only visible INSIDE of sourcefile.py, and not out here

print 'In \'parent process\', imported_source.c has value: ' + str(imported_source.c)
print 'Running lambda with protection'
print '-- begin lambda --'
run_lambda.executeFileWithProtection('sourcefile.py')
print '-- end lambda --'

print 'In \'parent process\', imported_source.c has value: ' + str(imported_source.c)
print ''

print 'Running lambda in the unprotected way'
print '-- begin lambda --'
run_lambda.executeFileWithoutProtection('sourcefile.py')
print '-- end lambda --'
print 'In \'parent process\', imported_source.c has value: ' + str(imported_source.c)


print '\n\nNext example\n\n'


print 'In \'parent process\', calling numpy arange: ' + str(numpy.arange(8))
print 'Run lambda that overwrites it (running in protected way)'
print '-- begin lambda --'
run_lambda.executeFileWithProtection('sourcefile_using_np.py')
print '-- finish lambda --'

print 'In \'parent process\', calling numpy arange again: ' + str(numpy.arange(8))
print ''

print 'Run lambda that overwrites it (run in unprotected way)'
print '-- begin lambda --'
run_lambda.executeFileWithoutProtection('sourcefile_using_np.py')
print '-- finish lambda --'
print 'In \'parent process\', calling numpy arange again: ' + str(numpy.arange)

