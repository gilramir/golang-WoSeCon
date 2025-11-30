package wosecon

var yesul_solution [][]rune = [][]rune{
        []rune{' ','마',' ',' ','바',' ',' ',' '}, // row 0
        []rune{' ',' ','트',' ','전','나',' ',' '}, // row 1
        []rune{'핸',' ',' ','비',' ',' ','나',' '}, // row 2
        []rune{'어','드','레',' ',' ','점',' ','지'}, // row 3
        []rune{'머','텔','폰',' ','화',' ','버',' '}, // row 4
        []rune{'니',' ',' ','백',' ','아','친','커'}, // row 5
        []rune{' ',' ',' ',' ',' ',' ','구','피'}, // row 6
        []rune{' ',' ','물','크','리','스','마','스'}, // row 7
}

var yesul_puzzle [][]rune = [][]rune{
        []rune{'통','마','요','얼','바','경','알','문'},
        []rune{'적','마','트','방','전','나','관','트'},
        []rune{'핸','청','민','비','키','김','나','참'},
        []rune{'어','드','레','한','의','점','한','지'},
        []rune{'머','텔','폰','신','화','의','버','담'},
        []rune{'니','집','결','백','다','아','친','커'},
        []rune{'수','한','타','다','에','에','구','피'},
        []rune{'피','원','물','크','리','스','마','스'},
}

var yesul_word_placements map[string]WordPlacement = map[string]WordPlacement{
    "바나나" :  WordPlacement{Row:0, Col: 4, Direction: LTRDescending},
    "아버지" :  WordPlacement{Row:5, Col: 5, Direction: LTRAscending},
    "어머니" :  WordPlacement{Row:3, Col: 0, Direction: Down},
    "친구":     WordPlacement{Row:5, Col: 6, Direction: Down},
    "백화점":   WordPlacement{Row:5, Col: 3, Direction: LTRAscending},
    "마트":     WordPlacement{Row:0, Col: 1, Direction: LTRDescending},
    "핸드폰":   WordPlacement{Row:2, Col: 0, Direction: LTRDescending},
    "크리스마스": WordPlacement{Row:7, Col: 3, Direction: LTRHorizontal},
    "텔레비전": WordPlacement{Row:4, Col:1, Direction: LTRAscending},
    "물":       WordPlacement{Row:7, Col:2, Direction: NilDirection},
    "커피":     WordPlacement{Row:5, Col:7, Direction: Down},
}

var yesul_banana =      "바나나"
var yesul_father =      "아버지"
var yesul_mother =      "어머니"
var yesul_friend =      "친구"
var yesul_deptstore =   "백화점"
var yesul_mart =        "마트"
var yesul_handphone=    "핸드폰"
var yesul_christmas =   "크리스마스"
var yesul_tv =          "텔레비전"
var yesul_water =       "물"
var yesul_coffee =      "커피"

var gangnam_style_vocab []string = []string{
    "오빤",
    "강남",
    "스타일",
    "사랑",
    "낮에는",
    "따사로운",
    "인간적인",
    "여자",
}

var ko_country_names []string = []string{
    "한국",
    "일본",
    "중국",
    "인도",
    "미국",
    "캐나다",
    "멕시코",
    "브라질",
    "이집트",
    "프랑스",
    "독일",
    "영국",
    "스페인",
    "이탈리아",
   // "태국",
   // "베트남",
   // "러시아",
   // "뉴질랜드",
   // "호주",
}

var ko_country_names_long []string = []string{
    "한국",
    "일본",
    "중국",
    "인도",
    "미국",
    "캐나다",
    "멕시코",
    "브라질",
    "이집트",
    "프랑스",
    "독일",
    "영국",
    "스페인",
    "이탈리아",
    "태국",
    "베트남",
    "러시아",
    "뉴질랜드",
    "호주",
}

var ko_kdh []string = []string{
    "귀신", "혼문", "사신", "악마",
    "봉인", "무대", "응원법", "마음", "자유",
}
