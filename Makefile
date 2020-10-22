
test1:
	./runtest.py --debug=test tests/step1_read_print.mal ./basilisk/parser.py

test2:
	./runtest.py --debug=test tests/step2_eval.mal ./basilisk/repl.py

test3:
	./runtest.py --debug=test tests/step3_env.mal ./basilisk/repl.py

test4:
	./runtest.py --debug=test tests/step4_if_fn_do.mal ./basilisk/repl.py

test5:
	./runtest.py --debug=test tests/step5_tco.mal ./basilisk/repl.py

test6:
	(cd ./basilisk && ../runtest.py --debug=test ../tests/step6_file.mal ./repl.py)

test7:
	(cd ./basilisk && ../runtest.py --debug=test ../tests/step7_quote.mal ./repl.py)

test8:
	(cd ./basilisk && ../runtest.py --debug=test ../tests/step8_macros.mal ./repl.py)

test9:
	(cd ./basilisk && ../runtest.py --debug=test ../tests/step9_try.mal ./repl.py)

testA:
	(cd ./basilisk && ../runtest.py --debug=test ../tests/stepA_mal.mal ./repl.py)
