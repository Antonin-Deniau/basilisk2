;; MACRO UTILS
(defmacro! default (fn* [arg def] `(if (nil? ~arg) ~def ~arg)))
(defmacro! defunc (fn* [name args func] `(def! ~name (fn* ~args ~func))))

;;UTILS
(defunc ord-between [low val max] (&& (>= (ord val) low) (<= (ord val) max)))

(defunc or-list [& opt] 
	(if (first opt)
	  true
	  (if (empty? opt) 
	    false
	    (apply or-list (rest opt)))))

;; MATCHER
(defunc number-matcher [c] (or-list (= c ".") (= c "-") (ord-between 48 c 57)))
(defunc vector-matcher [c] (= "[" c))
(defunc list-matcher [c] (= "(" c))
(defunc map-matcher [c] (= "{" c))
(defunc keyword-matcher [c] (= ":" c))
(defunc whitespace-matcher [c] (or-list (= " " c) (= "\n" c) (= "\t" c)))


;; READER
(defunc whitespace-ignore [stream] 
	(if (whitespace-matcher (peek-byte stream))
	  (do
	    (read-byte stream)
	    (whitespace-ignore stream))
	  stream))


(defunc number-reader-iterate-decimals [index positive res reader-macro stream]
	(let* [c (peek-byte stream)]
	  (if (= c "")
	    (if positive res (* res -1))
	    (if (ord-between 48 c 57)
	      (number-reader-iterate-decimals (+ index 1) positive (+ res (/ (- (ord (read-byte stream)) 48) (** 10 index))) reader-macro stream)
	      (if positive res (* res -1))))))

(defunc number-reader-iterate [positive res reader-macro stream]
	(let* [c (peek-byte stream)]
	  (if (= c "")
	    (if positive res (* res -1))
	    (if (= c ".")
	      (do
		(read-byte stream)
		(number-reader-iterate-decimals 1 positive res reader-macro stream))
	      (if (ord-between 48 c 57)
		(number-reader-iterate positive (+ (* res 10) (- (ord (read-byte stream)) 48)) reader-macro stream)
		(if positive res (* res -1)))))))

(defunc number-reader [reader-macro stream] 
	(let* [negative (= (peek-byte stream) "-")]
	  (if negative
	    (do
	      (read-byte stream)
	      (number-reader-iterate false 0 reader-macro stream))
	    (number-reader-iterate true 0 reader-macro stream))))

(defunc vector-reader-iterate [ret reader-macro stream]
	(let* [stream (whitespace-ignore stream)] ;; IGNORE WHITESPACE
	  (if (= (peek-byte stream) "]") ;; RETURN ON VECTOR END
	    (do
	      (read-byte stream)
	      ret)
	    (let* [data (read reader-macro stream)]
	      (vector-reader-iterate (conj ret data) reader-macro stream)))))

(defunc vector-reader [reader-macro stream] 
	(do
	  (read-byte stream)
	  (vector-reader-iterate [] reader-macro stream)))

(defunc list-reader-iterate [ret reader-macro stream]
	(let* [stream (whitespace-ignore stream)] ;; IGNORE WHITESPACE
	  (if (= (peek-byte stream) ")") ;; RETURN ON VECTOR END
	    (do
	      (read-byte stream)
	      ret)
	    (let* [data (read reader-macro stream)]
	      (list-reader-iterate (conj ret data) reader-macro stream)))))

(defunc keyword-reader-iterate [res reader-macro stream]
	(let* [c (peek-byte stream)]
	  (if (= c "")
	    (if (= res "")
	      (raise "Unable to get keyword.")
	      (keyword res))
	    (if (or-list (ord-between 65 c 90) (= c "-") (= c "_") (ord-between  96 c 122))
	      (do
		(keyword-reader-iterate (str res (read-byte stream)) reader-macro stream))
	      (keyword res)))))

(defunc keyword-reader [reader-macro stream]
	(do
	  (read-byte stream)
	  (keyword-reader-iterate "" reader-macro stream)))

(defunc list-reader [reader-macro stream] 
	(do
	  (read-byte stream)
	  (list-reader-iterate () reader-macro stream)))

(defunc map-reader-iterate [ret reader-macro stream]
	(let* [stream (whitespace-ignore stream)] ;; IGNORE WHITESPACE
	  (if (= (peek-byte stream) "}") ;; RETURN ON VECTOR END
	    (do
	      (read-byte stream)
	      ret)
	    (let* [data (read reader-macro stream)]
	      (map-reader-iterate (conj ret data) reader-macro stream)))))

(defunc map-reader [reader-macro stream] 
	(do
	  (read-byte stream)
	  (let* [items (map-reader-iterate [] reader-macro stream)]
	    (if (= (% (count items) 2) 0)
	      (apply hash-map items)
	      (raise "Items in hash-map must be in pair")))))


;; READER MACRO
(def! reader-macro [[keyword-matcher keyword-reader]
		    [map-matcher     map-reader]
		    [list-matcher    list-reader]
		    [number-matcher  number-reader]
		    [vector-matcher  vector-reader]])

;; PARSER FUNCTION
(defunc match [reader-macro c]
	(if (empty? reader-macro)
	  nil
	  (let* [macro (first reader-macro)
		       matcher (nth macro 0)
		       reader  (nth macro 1)
		       is-true (matcher c)]
	    (if is-true
	      reader
	      (match (rest reader-macro) c)))))

(defunc read [reader-macro stream]
	(let* [stream (whitespace-ignore stream)
		      byte   (peek-byte stream)
		      reader (match reader-macro byte)]
	  (if (nil? reader)
	    (raise (str "Unable to find matcher: [" byte "]"))
	    (reader reader-macro stream))))

;; BASIC PARSER SETUP
(prn (read reader-macro (string-stream "(1 [-1.5554 .2 5484.263 { :a 5 :s-_ 5 :s-_a 8} 5 -1] 5)")))
