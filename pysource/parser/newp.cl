;; MACRO UTILS
(defmacro! defun (fn* [name args func] `(def! ~name (fn* ~args ~func))))

;;UTILS
(defun ord-between [low val max] (&& (>= (ord val) low) (<= (ord val) max)))

(defun or-list [& opt] 
  (if (first opt)
    true
    (if (empty? opt) 
      false
      (apply or-list (rest opt)))))

;; MATCHER
(defun number-matcher [c] (or-list (= c ".") (= c "-") (ord-between 48 c 57)))
(defun vector-matcher [c] (= "[" c))
(defun list-matcher [c] (= "(" c))
(defun map-matcher [c] (= "{" c))
(defun keyword-matcher [c] (= ":" c))
(defun string-matcher [c] (= "\"" c))
(defun quote-matcher))
(defun deref-matcher [c] (= "@" c))
(defun unquote-matcher [c] (= "~" c))
(defun quasiquote-matcher [c] (= "`" c))
(defun metadata-matcher [c] (= "^" c))
(defun whitespace-matcher [c] (or-list (= " " c) (= "\n" c) (= 9 (ord c))))
(defun symbol-matcher [c] (! (whitespace-matcher c)))


;; READER
(defun ignore-until-newline [stream]
  (if (= (read-byte stream) "\n")
    stream
    (ignore-until-newline stream)))

(defun whitespace-ignore [stream] 
  (if (whitespace-matcher (peek-byte stream))
    (do
      (read-byte stream)
      (whitespace-ignore stream))
    (if (= (peek-byte stream) ";")
      (whitespace-ignore (ignore-until-newline stream))
      stream)))

(defun number-reader-iterate-decimals [index positive res reader-macro stream]
  (let* [c (peek-byte stream)]
    (if (= c "")
      (if positive res (* res -1))
      (if (ord-between 48 c 57)
        (number-reader-iterate-decimals (+ index 1) positive (+ res (/ (- (ord (read-byte stream)) 48) (** 10 index))) reader-macro stream)
        (if positive res (* res -1))))))

(defun number-reader-iterate [positive res reader-macro stream]
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


(defun number-reader [reader-macro stream] 
  (let* [negative (= (peek-byte stream) "-")]
    (if negative
      (do
        (read-byte stream)
        (number-reader-iterate false 0 reader-macro stream))
      (number-reader-iterate true 0 reader-macro stream))))

(defun deref-reader [reader-macro stream]
  (do
    (read-byte stream)
    (list 'deref (read reader-macro stream))))

(defun quote-reader [reader-macro stream]
  (do
    (read-byte stream)
    (list 'quote (read reader-macro stream))))

(defun metadata-reader [reader-macro stream]
  (do
    (read-byte stream)
    (list 'with-meta (read reader-macro stream) (read reader-macro stream))))

(defun unquote-reader [reader-macro stream]
  (do
    (read-byte stream)
    (if (= (peek-byte stream) "@")
      (do
        (read-byte stream)
        (list 'splice-unquote (read reader-macro stream)))
      (list 'unquote (read reader-macro stream)))))

(defun quasiquote-reader [reader-macro stream]
  (do
    (read-byte stream)
    (list 'quasiquote (read reader-macro stream))))

(defun vector-reader-iterate [ret reader-macro stream]
  (let* [stream (whitespace-ignore stream)]
    (if (= (peek-byte stream) "]")
      (do
        (read-byte stream)
        ret)
      (let* [data (read reader-macro stream)]
        (vector-reader-iterate (conj ret data) reader-macro stream)))))

(defun vector-reader [reader-macro stream] 
  (do
    (read-byte stream)
    (vector-reader-iterate [] reader-macro stream)))

(defun list-reader-iterate [ret reader-macro stream]
  (let* [stream (whitespace-ignore stream)]
    (if (= (peek-byte stream) ")")
      (do
        (read-byte stream)
        ret)
      (let* [data (read reader-macro stream)]
        (list-reader-iterate (conj ret data) reader-macro stream)))))

(defun keyword-reader-iterate [res reader-macro stream]
  (let* [c (peek-byte stream)]
    (if (= c "")
      (if (= res "")
        (raise "Unable to get keyword.")
        (keyword res))
      (if (or-list (ord-between 65 c 90) (= c "-") (= c "_") (ord-between  96 c 122))
        (keyword-reader-iterate (str res (read-byte stream)) reader-macro stream)
        (keyword res)))))

(defun keyword-reader [reader-macro stream]
  (do
    (read-byte stream)
    (keyword-reader-iterate "" reader-macro stream)))

(defun symbol-reader-iterate [res reader-macro stream]
  (let* [c (peek-byte stream)]
    (if (= c "")
      (if (= res "")
        (raise "Unable to get symbol.")
        res)
      (if (or-list (= " " c) (= "\n" c) (= 9 (ord c)))
        res
        (symbol-reader-iterate (str res (read-byte stream)) reader-macro stream)))))

(defun filter-special-symbols [s]
  (cond
    (= s "true") true
    (= s "false") false
    (= s "nil") nil
    "else" (symbol s)))

(defun symbol-reader [reader-macro stream]
  (filter-special-symbols (symbol-reader-iterate "" reader-macro stream)))

(defun list-reader [reader-macro stream] 
  (do
    (read-byte stream)
    (list-reader-iterate () reader-macro stream)))

(defun map-reader-iterate [ret reader-macro stream]
  (let* [stream (whitespace-ignore stream)]
    (if (= (peek-byte stream) "}")
      (do
        (read-byte stream)
        ret)
      (let* [data (read reader-macro stream)]
        (map-reader-iterate (conj ret data) reader-macro stream)))))

(defun map-reader [reader-macro stream] 
  (do
    (read-byte stream)
    (let* [items (map-reader-iterate [] reader-macro stream)]
      (if (= (% (count items) 2) 0)
        (apply hash-map items)
        (raise "Items in hash-map must be in pair")))))


(defun string-reader-iterate [res esc reader-macro stream]
  (let* [c (read-byte stream)]
    (cond
      (= c "") (raise "Unexpected end of input")
      (= c "n") (if esc
                    (string-reader-iterate (str res "\n") false reader-macro stream)
                    (string-reader-iterate (str res c) false reader-macro stream))
      (= c "\"") (if esc 
                    (string-reader-iterate (str res c) false reader-macro stream)
                    res)
      (= c "\\") (if esc
                    (string-reader-iterate (str res c) false reader-macro stream)
                    (string-reader-iterate res true reader-macro stream))
      "else" (string-reader-iterate (str res c) false reader-macro stream))))

(defun string-reader [reader-macro stream]
  (do 
    (read-byte stream)
    (string-reader-iterate "" false reader-macro stream)))

;; READER MACRO
(def! reader-macro [[keyword-matcher    keyword-reader]
                    [map-matcher        map-reader]
                    [list-matcher       list-reader]
                    [number-matcher     number-reader]
                    [vector-matcher     vector-reader]
                    [string-matcher     string-reader]
                    [quote-matcher      quote-reader]
                    [unquote-matcher    unquote-matcher]
                    [quasiquote-matcher quasiquote-matcher]
                    [quote-matcher      quote-reader]
                    [metadata-matcher   metadata-reader]
                    [deref-matcher      deref-reader]
                    [symbol-matcher     symbol-reader]])

;; PARSER FUNCTION
(defun match [reader-macro c]
  (if (empty? reader-macro)
    nil
    (let* [macro   (first reader-macro)
           matcher (nth macro 0)
           reader  (nth macro 1)
           is-true (matcher c)]
      (if is-true
        reader
        (match (rest reader-macro) c)))))

(defun read [reader-macro stream]
  (let* [stream (whitespace-ignore stream)
         byte   (peek-byte stream)
         reader (match reader-macro byte)]
    (if (nil? reader)
      (raise (str "Unable to find matcher for byte: `" byte "`"))
      (reader reader-macro stream))))

;; BASIC PARSER SETUP
;(prn (read reader-macro (string-stream "(1 [-1.5554 name \"\\\"loli\\nl\\rol\" nil true 'false '(1 2) .2 5484.263 { :a 5 :s-_ 5 :s-_a 8} 5 -1] 5)")))
(def! file (output-stream "./parser/newp.cl"))
(prn (read reader-macro file))
