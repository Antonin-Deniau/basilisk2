test1:
	./runtest.py --debug=test tests/step1_read_print.mal ./src/parser.py

test2:
	./runtest.py --debug=test tests/step2_eval.mal ./src/basilisk

test3:
	./runtest.py --debug=test tests/step3_env.mal ./src/basilisk

test4:
	./runtest.py --debug=test tests/step4_if_fn_do.mal ./src/basilisk

test5:
	./runtest.py --debug=test tests/step5_tco.mal ./src/basilisk

test6:
	(cd ./src && ../runtest.py --debug=test ../tests/step6_file.mal ./basilisk)

test7:
	(cd ./src && ../runtest.py --debug=test ../tests/step7_quote.mal ./basilisk)

test8:
	(cd ./src && ../runtest.py --debug=test ../tests/step8_macros.mal ./basilisk)

test9:
	(cd ./src && ../runtest.py --debug=test ../tests/step9_try.mal ./basilisk)

testA:
	(cd ./src && ../runtest.py --debug=test ../tests/stepA_mal.mal ./basilisk)

testA1:
	(cd ./mal && ../runtest.py --debug=test ../tests/step1_read_print.mal basilisk ./step1_read_print.mal)

testA2:
	(cd ./mal && ../runtest.py --debug=test ../tests/step2_eval.mal basilisk ./step2_eval.mal)

testA3:
	(cd ./mal && ../runtest.py --debug=test ../tests/step3_env.mal basilisk ./step3_env.mal)
