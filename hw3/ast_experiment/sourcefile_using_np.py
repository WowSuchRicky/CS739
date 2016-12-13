import numpy

print 'Using numpy.arange: ' +  str(numpy.arange(5))

value_to_change_to = 10
print 'Changing value of numpy.arange to an integer: ' + str(value_to_change_to)
numpy.arange = value_to_change_to

print 'Using numpy.arange again: ' + str(numpy.arange)
