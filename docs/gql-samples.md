# Sample GQL Queries and Mutations

## Vocab

query FindVocab 
  {
  vocab(id: 2856) {
    learning_lang
    first_lang
    alternatives
    skill
    infinitive
    pos
    hint
    num_learning_words
    known_lang_code
    learning_lang_code
  }
}

query FindVocabs {
  vocabs(learning_code: "es",
  has_first: false,
  limit: 4) {
    id
    learning_lang
    first_lang
    alternatives
    skill
    infinitive
    pos
    hint
    num_learning_words
    known_lang_code
    learning_lang_code
  }
}

mutation UpdateVocab {
  updateVocab(input: {
    id: "1865",
    first_lang: "Mexico",
    alternatives: "",
    skill: "Travel",
    infinitive: "",
    pos: "Proper noun",
    hint: "",
    num_learning_words: 1
  }) {
    id
    learning_lang
    first_lang
    alternatives
    skill
    infinitive
    pos
    hint
    num_learning_words
    known_lang_code
    learning_lang_code
  }
}

mutation CreateVocab {
  createVocab(input: {
    learning_lang: "plateado",
    first_lang: "silver",
    alternatives: "",
    skill: "Colors",
    infinitive: "",
    pos: "adjective",
    hint: "you could plate something with it",
    num_learning_words: 1,
    known_lang_code: "en",
    learning_lang_code: "es",
  }) {
    id
    learning_lang
    first_lang
    alternatives
    skill
    infinitive
    pos
    hint
    num_learning_words
    known_lang_code
    learning_lang_code
  }
}


# Audits

query FindAudits {
  audits(
    table_name: "fixit",
    object_id : 0,
    start_time: "",
    end_time: "",
    limit: 40
  ) {
    id
    table_name
    object_id
    before
    after
    diff
    comments
    created
    created_by
  }
}

query FindAudit {
  audit(
    id:2
  ) {
    id
    table_name
    object_id
    before
    after
    diff
    comments
    created
    created_by
  }
}

# Fixits

mutation CreateFixit {
  createFixit(input: {
    vocab_id: "2799",
    status: PENDING,
    comments: "user asked for this to get fixed",
    field_name: "learning_lang"
  }) {
  id
  }
}

query FindFixit {
  fixit(id: 1) {
    id
    field_name
    created
    comments
    vocab_id
    status
    created_by
  }
}

query FindFixits {
  fixits(
    status: COMPLETED,
    vocab_id: 0,
    start_time:"",
    end_time: "",
    limit: 20
  ) {
    id
    field_name
    created
    comments
    vocab_id
    status
    created_by
  }
}

mutation updateFixit {
  updateFixit(input: {
    id: 2,
    status: PENDING,
    comments: "user request",
    field_name: "learning_lang"
  }) {
    id
    vocab_id
    status
    comments
    field_name
  }
}
