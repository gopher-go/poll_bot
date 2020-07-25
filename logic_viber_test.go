package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWeHaveNewMessageAfterUnsubscribe(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()
	userID := "124"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Equal(t, "Добрый день!\nВы гражданин Республики Беларусь?", reply.text)
	require.Equal(t, []string{"1. Да", "2. Нет"}, reply.options)

	text := newTextCallback(t, userID, "1. Да")
	require.Equal(t, text.User.ID, userID)
	reply, err = generateReplyFor(p, s, text)
	require.NoError(t, err)
	require.Equal(t, "Укажите, пожалуйста, Ваш возраст", reply.text)

	reply, err = generateReplyFor(p, s, newUnsubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Nil(t, reply)

	reply, err = generateReplyFor(p, s, newSubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Equal(t, "Добрый день!\nУкажите, пожалуйста, Ваш возраст", reply.text)
	require.Equal(t, []string{"1. младше 18", "2. 18-25", "3. 26-40", "4. 41-55", "5. старше 55"}, reply.options)
}

func TestUserFlowCaseSensitive(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	userID := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Equal(t, "Добрый день!\nВы гражданин Республики Беларусь?", reply.text)
	require.Equal(t, []string{"1. Да", "2. Нет"}, reply.options)

	text := newTextCallback(t, userID, "Привет")
	require.Equal(t, text.User.ID, userID)
	reply, err = generateReplyFor(p, s, text)
	require.NoError(t, err)
	require.Equal(t, "Пожалуйста выберите один из предложенных ответов или введите его номер.\n\nВы гражданин Республики Беларусь?", reply.text)
	require.Equal(t, []string{"1. Да", "2. Нет"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "1. Да"))
	require.NoError(t, err)
	require.Equal(t, "Укажите, пожалуйста, Ваш возраст", reply.text)
	require.Equal(t, []string{"1. младше 18", "2. 18-25", "3. 26-40", "4. 41-55", "5. старше 55"}, reply.options)

	user, err := s.fromPersisted(userID)
	require.NoError(t, err)
	require.Equal(t, user.ID, userID)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "i'm a tester, clear it"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Your storage cleared with <nil>")

	user, err = s.fromPersisted(userID)
	require.NoError(t, err)
	require.Nil(t, user)
}

func TestCaseInsensitive(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	caseSensitive = false
	userID := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Добрый день!\nВы гражданин Республики Беларусь?")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста выберите один из предложенных ответов или введите его номер.\n\nВы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "2. нету"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста выберите один из предложенных ответов или введите его номер.\n\nВы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "1. ДА"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "1. младше 18"))
	require.NoError(t, err)
	require.Equal(t, "Вам должно быть 18 или больше.\n\nУкажите, пожалуйста, Ваш возраст", reply.text)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "4. 41-55"))
	require.NoError(t, err)
	require.Equal(t, "К какому типу относится ваш населенный пункт?", reply.text)
	require.Equal(t, reply.options, []string{"1. Областной центр", "2. Город или городской поселок", "3. Агрогородок / Село / Деревня", "4. Проживаю за пределами РБ"})
	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "3. Агрогородок / село / Деревня"))
	require.NoError(t, err)

	require.Equal(t, "За кого Вы планируете проголосовать?", reply.text)
	require.Equal(t, []string{"1. Дмитриев", "2. Канопацкая", "3. Лукашенко", "4. Тихановская", "5. Черечень", "6. Против всех", "7. Затрудняюсь ответить", "8. Не пойду голосовать"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "1. Дмитриев"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш пол")
	require.Equal(t, reply.options, []string{"1. Мужской", "2. Женский"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "1. мужской"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Выберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск")
	require.Equal(t, reply.options, []string{"1. Брестская", "2. Витебская", "3. Гомельская", "4. Гродненская", "5. Минская", "6. Могилевская", "7. Минск", "8. Проживаю за пределами РБ"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "2. витебская"))
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, "Ваш уровень образования?", reply.text)
	require.Equal(t, []string{"1. Базовое / Среднее общее (школа)", "2. Профессионально-техническое", "3. Среднее специальное", "4. Высшее", "5. Другое"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "3. Среднее специальное"))
	require.NoError(t, err)
	require.Equal(t, "Укажите, пожалуйста, общий совокупный месячный доход вашей семьи (включая пенсии, стипендии, надбавки и прочее)", reply.text)
	require.Equal(t, reply.options, []string{"1. До 500 бел. руб.", "2. 500 - 1000 бел. руб.", "3. 1000 - 2000 бел. руб.", "4. Выше 2000 бел.руб.", "5. Не хочу отвечать на этот вопрос"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "1. До 500 бел. руб."))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Спасибо, ваш голос учтен!")

	user, err := s.fromPersisted(userID)
	require.NoError(t, err)
	require.Equal(t, user.ID, userID)
	require.Equal(t, user.Properties["age"], "4. 41-55")
	require.Equal(t, 8, user.Level)
	require.Equal(t, user.Candidate, "1. Дмитриев")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Спасибо, ваш голос учтен!")

	reply, err = generateReplyFor(p, s, newUnsubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Nil(t, reply)

	reply, err = generateReplyFor(p, s, newSubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Equal(t, "Добрый день!\nВы уже приняли участие в Народном опросе. Спасибо, ваш голос учтен!", reply.text)
}

func TestUserFlow(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	userID := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Equal(t, "Добрый день!\nВы гражданин Республики Беларусь?", reply.text)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста выберите один из предложенных ответов или введите его номер.\n\nВы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "1. Да"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "1. младше 18"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Вам должно быть 18 или больше.\n\nУкажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "4. 41-55"))
	require.NoError(t, err)

	require.Equal(t, "К какому типу относится ваш населенный пункт?", reply.text)

	user, err := s.fromPersisted(userID)
	require.NoError(t, err)

	require.Equal(t, user.ID, userID)
	fmt.Println(user)
	require.Equal(t, user.Properties["age"], "4. 41-55")
	require.Equal(t, 2, user.Level)

	seenCallback := newSeenCallback(t, userID)
	require.Equal(t, seenCallback.User.ID, userID)
	reply, err = generateReplyFor(p, s, seenCallback)
	require.NoError(t, err)
	require.Nil(t, reply)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "3. Агрогородок / Село / Деревня"))
	require.NoError(t, err)
	require.Equal(t, "За кого Вы планируете проголосовать?", reply.text)
	require.Equal(t, []string{"1. Дмитриев", "2. Канопацкая", "3. Лукашенко", "4. Тихановская", "5. Черечень", "6. Против всех", "7. Затрудняюсь ответить", "8. Не пойду голосовать"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "3. Лукашенко"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш пол")
	require.Equal(t, reply.options, []string{"1. Мужской", "2. Женский"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "1. Мужской"))
	require.NoError(t, err)

	user, err = s.fromPersisted(userID)
	require.NoError(t, err)

	require.Equal(t, user.ID, userID)
	require.Equal(t, user.Properties["age"], "4. 41-55")
	require.Equal(t, 5, user.Level)
	require.Equal(t, "3. Лукашенко", user.Candidate)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, "Пожалуйста выберите один из предложенных ответов или введите его номер.\n\nВыберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск", reply.text)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, "Пожалуйста выберите один из предложенных ответов или введите его номер.\n\nВыберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск", reply.text)

	subscribe := newSubscribeCallback(t, userID)
	user, err = s.Obtain(userID)
	require.NoError(t, err)
	reply, err = generateReplyFor(p, s, subscribe)
	require.NoError(t, err)
	require.Empty(t, reply)

	reply, err = generateReplyFor(p, s, newSeenCallback(t, userID))
	require.NoError(t, err)
	require.Empty(t, reply)
}

func TestNumericReplies(t *testing.T) {
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)

	p := generateOurPoll()

	userID := "123"

	reply, err := generateReplyFor(p, s, newSubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Equal(t, "Добрый день!\nВы гражданин Республики Беларусь?", reply.text)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста выберите один из предложенных ответов или введите его номер.\n\nВы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"1. Да", "2. Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "2"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "4"))
	require.NoError(t, err)
	require.Equal(t, "К какому типу относится ваш населенный пункт?", reply.text)

	user, err := s.fromPersisted(userID)
	require.NoError(t, err)

	require.Equal(t, user.ID, userID)
	require.Equal(t, 2, user.Level)
	require.Equal(t, user.Properties["age"], "4. 41-55")
}

func newUnsubscribeCallback(t *testing.T, id string) *ViberCallback {
	json := `{"event":"unsubscribed","timestamp":1595347885535,"chat_hostname":"SN-376_","user_id":"%s","message_token":5466394919049723652}`
	validJSON := fmt.Sprintf(json, id)

	ret, err := parseCallback([]byte(validJSON))
	require.NoError(t, err)

	return ret
}

func newSubscribeCallback(t *testing.T, id string) *ViberCallback {
	c := &ViberCallback{
		Event: "subscribed",
		User: User{
			ID:   id,
			Name: "Vasya",
		},
	}

	b, err := json.Marshal(c)
	require.NoError(t, err)

	ret, err := parseCallback(b)
	require.NoError(t, err)

	return ret
}

func newTextCallback(t *testing.T, id string, text string) *ViberCallback {
	json := `{"event":"message","sender":{"id":"%s","Name":"Vasya"},"message":{"type":"text","text":"%s"}}`

	validJSON := fmt.Sprintf(json, id, text)

	ret, err := parseCallback([]byte(validJSON))
	require.NoError(t, err)

	return ret
}

func newSeenCallback(t *testing.T, id string) *ViberCallback {
	json := `{"event":"seen","user_id":"%s"}`

	validJSON := fmt.Sprintf(json, id)

	ret, err := parseCallback([]byte(validJSON))
	require.NoError(t, err)

	return ret
}
