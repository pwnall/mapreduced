import glob
import os.path

f = open('output', 'w')

golden_files = glob.glob('cases/*.out')
for golden_file in sorted(golden_files):
    test_name = os.path.basename(golden_file)
    output_file = '/usr/testguest/' + test_name
    if not os.path.isfile(output_file):
        f.write(test_name + ' not produced by test program\n')
        continue

    with open(golden_file, 'r') as g:
        golden_data = g.read()
    with open(output_file, 'r') as o:
        output_data = o.read()
    if golden_data == output_data:
        f.write(test_name + ' matches expected output\n')
    else:
        f.write(test_name + ' does not match expected output\n')

f.close()
