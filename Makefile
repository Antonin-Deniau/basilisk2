gobuild:
	go run parser.go types.go

gotest:
	./runtest.py --debug=test tests/step2_eval.mal ./basilisk

test2:
	./runtest.py --debug=test tests/step2_eval.mal ./basilisk

test3:
	./runtest.py --debug=test tests/step3_env.mal ./basilisk

test4:
	./runtest.py --debug=test tests/step4_if_fn_do.mal ./basilisk

test5:
	./runtest.py --debug=test tests/step5_tco.mal ./basilisk

test6:
	./runtest.py --debug=test ../tests/step6_file.mal ./basilisk

test7:
	./runtest.py --debug=test ../tests/step7_quote.mal ./basilisk

test8:
	./runtest.py --debug=test ../tests/step8_macros.mal ./basilisk

test9:
	./runtest.py --debug=test ../tests/step9_try.mal ./basilisk

testA:
	./runtest.py --debug=test ../tests/stepA_mal.mal ./basilisk

testA1:
	(cd ./mal && ../runtest.py --debug=test ../tests/step1_read_print.mal ../basilisk ./step1_read_print.mal)

testA2:
	(cd ./mal && ../runtest.py --debug=test ../tests/step2_eval.mal ../basilisk ./step2_eval.mal)

testA3:
	(cd ./mal && ../runtest.py --debug=test ../tests/step3_env.mal ../basilisk ./step3_env.mal)

testA4:
	(cd ./mal && ../runtest.py --debug=test ../tests/step4_if_fn_do.mal ../basilisk ./step4_if_fn_do.mal)

testA6:
	(cd ./mal && ../runtest.py --debug=test ../tests/step6_file.mal ../basilisk ./step6_file.mal)

testA7:
	(cd ./mal && ../runtest.py --debug=test ../tests/step7_quote.mal ../basilisk ./step7_quote.mal)

testA8:
	(cd ./mal && ../runtest.py --debug=test ../tests/step8_macros.mal ../basilisk ./step8_macros.mal)

testA9:
	(cd ./mal && ../runtest.py --debug=test ../tests/step9_try.mal ../basilisk ./step9_try.mal)

