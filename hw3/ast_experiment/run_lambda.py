import ast

'''
The idea: 
We have a main process with some globally imported modules, and we want to be able to import another 
module and run the code within that module, but prevent that newly imported module from modifying the 
globally imported modules. The idea is to pick out any writes to those globally imported modules, and
change the namespace to something else. We can do some static analysis to handle this

Process to do that:
// IGNORE 1) maintain a list of globally imported modules

-- begin lambda import --
2) build a list of things imported in the lambda
// IGNORE 3) get the intersection of the two lists
   - if empty, nothing left to do so stop
   - if non-empty, continue
4) (pass #1) visit each assignment (i.e. WRITE!) statement in the AST
   if left hand side is in our list of namespaces, then note the namespace and
   variable name
5) once we've done that entire pass, we have a list of variables in separate namespaces
   that are being written to. These are essentially what we need to change.
6) For each namespace in that list, create a new uniquely named namespace object (use 
   a dictionary to keep track of these; key is old namespace name, value is new one). We
   need to do this in the actual code, so modify the AST in this case
7) (think of namespace as a dict) - for each variable in that list, add an entry in 
   the namespace where the key is the variable name and the value is the reference to
   the older namespace (i.e. we would make a new namespace for numpy and if the
   program has numpy.arange = 4, we add numpy.arange to our list so we make a new 
   namespace called numpy2 (or something like that), we add key arange to it with value
   numpy.arange.
6) (pass #2) visit each name usage from module and if its reerring to something that was in our variable list, change the namespace to the new namespace

-- end lambda import --
7) we now have some changed lambda code; we can store this somewhere so we don't need to go through this importation process repeatedly (makes sense if constantly executing similar lambda)

ASSUMPTIONS:
 (big deal, could be deal breaker): not handling recursive namespace differences
'''

# good resource: http://stackoverflow.com/questions/1515357/simple-example-of-how-to-use-ast-nodevisitor
# another, about namespaces: http://stackoverflow.com/questions/28345780/how-do-i-create-a-python-namespace-argparse-parse-args-value
# another about using string as keyword argument: http://stackoverflow.com/questions/2932648/how-do-i-use-a-string-as-a-keyword-argument
# take a look at this: http://stackoverflow.com/questions/37281928/making-a-copy-of-an-entire-namespace
# interesting read on exec: http://lucumr.pocoo.org/2011/2/1/exec-in-python/


DEBUG = False

# recurse down an Attribute node and return full string
def recurseAttribute(node):
    if isinstance(node, ast.Attribute):
        return recurseAttribute(node.value) + "." + node.attr;
    elif isinstance(node, ast.Name):
        return node.id
    else:
        return ''

def changeNamespace(node):
    if isinstance(node, ast.Attribute):
        changeNamespace(node.value)
    elif isinstance(node, ast.Name):
        node.id = new_namespaces[node.id]

imported_modules = []
assigned_variables = []
new_namespaces = {} # key is old namespace name, value is new namespace name
namespace_vars = {} # key is new namespace name, value is list of variables being written in it

class RemoveImports(ast.NodeTransformer):
    def visit_Import(self, node):
        if (DEBUG):
            print 'recording and removing import'
        mod_names = node.names
        for x in mod_names:
            imported_modules.append(x.name)
        return None


class FindAssignedVars(ast.NodeVisitor):
    def visit_Assign(self, node):
        if (DEBUG):
            print 'assignment visit'
        # go through targets; these are the vars being assigned to
        for x in node.targets:
            # only if namespace (i.e. Attribute)
            if isinstance(x, ast.Attribute):
                assigned_variables.append(recurseAttribute(x))

class ChangeRelevantNamespaces(ast.NodeTransformer):
    def visit_Attribute(self, node):
        var_delim = recurseAttribute(node).split('.')
        if var_delim[0] in new_namespaces: # if it's an old namespace
            if var_delim[1] in namespace_vars[new_namespaces[var_delim[0]]]:
                if (DEBUG):
                    print 'Changing namespace reference'
                changeNamespace(node)
        return node

def executeFileWithProtection(path):

    # get the ast of the lambda
    f = open(path, 'r')
    source = f.read()
    node = ast.parse(source, mode='exec')

    if (DEBUG):
        print ast.dump(node)

    # get list of imported modules, and list of variables being
    # assigned to that are in a different namespace
    RemoveImports().visit(node)
    FindAssignedVars().visit(node)

    if (DEBUG):
        print 'Imported modules: '
        print imported_modules
        print 'Assigned variables in different namespace: '
        print assigned_variables

    # go through the list of assigned variables and:
    # 1) pick out the actual namespace and decide on name of new namespace to use
    # 2) take note of the variable IN that namespace (we will have to copy this to the
    #    new namespace later
    for x in assigned_variables:
        old_namespace = x.split('.')[0]
        new_namespace = old_namespace + '2'
        var_name = x.split('.')[1]

        if old_namespace not in new_namespaces:
            new_namespaces[old_namespace] = new_namespace

        if new_namespace not in namespace_vars:
            namespace_vars[new_namespace] = []

        old_list = namespace_vars[new_namespace]
        old_list.append(var_name)

    if (DEBUG):
        print 'Dictionary of namespaces (key = old, value = new)'
        print new_namespaces
        print 'Dictionary of namespace to variables being written'
        print namespace_vars
    
    # now we need to modify the actual provided code/AST; we'll actually do this
    # by creating a new AST (maybe there's a more efficient way)
    # 1) add a namespace object for each new namespace (key in namespace_vars)
    # 2) add an assignment statement with:
    #    LHS: new_namespace.var_name
    #    RHS: old_namespace.var_name 

    codeToAdd = "class Namespace:\n"
    codeToAdd += "    def __init__(self, **kwargs):\n"
    codeToAdd += "        self.__dict__.update(kwargs)\n"
    for old_ns, new_ns in new_namespaces.iteritems():
        codeToAdd += (new_ns + " = Namespace()\n") 
        for var_name in namespace_vars[new_ns]:
            codeToAdd += (new_ns + "." + var_name + " = " + old_ns + "." + var_name + "\n")

    import_stmts = ""
    for x in imported_modules:
        import_stmts += 'import ' + x + '\n'

    codeToAdd = import_stmts + codeToAdd

    new_ast = ast.parse(codeToAdd, mode = "exec")

    # merge ASTs
    # NOTE: for it to be correct, new_code_ast must come after import statements, 
    #       which is why we pull out those import statements earlier
    new_ast.body += node.body

    # now that we've done that, go back through the AST and change each assignment
    # statement that was writing to a different namespace such that it is now writing
    # to the new namespace (update AST)

    if (DEBUG):
        print 'Second pass, changing assignments now'

    ChangeRelevantNamespaces().visit(node)

    if (DEBUG):
        print 'Updated ast is: '
        print ast.dump(node)

    # compile execute the code (updated code, if we changed AST)
    codeobj = compile(new_ast, filename="<ast>", mode="exec")
    exec(codeobj)


def executeFileWithoutProtection(path):
    f = open(path, 'r')
    source = f.read()
    node = ast.parse(source, mode='exec')
    codeobj = compile(node, filename="<ast>", mode="exec")
    exec(codeobj)
