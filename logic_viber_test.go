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
	require.Equal(t, fmt.Sprintf(welcomeHeader, 568)+"\n\nВы гражданин Республики Беларусь?", reply.text)
	require.Equal(t, []string{"Да", "Нет"}, reply.options)

	text := newTextCallback(t, userID, "Да")
	require.Equal(t, text.User.ID, userID)
	reply, err = generateReplyFor(p, s, text)
	require.NoError(t, err)
	require.Equal(t, "Укажите, пожалуйста, Ваш возраст", reply.text)

	reply, err = generateReplyFor(p, s, newUnsubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Nil(t, reply)

	reply, err = generateReplyFor(p, s, newSubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf(welcomeHeader, 569)+"\n\nУкажите, пожалуйста, Ваш возраст", reply.text)
	require.Equal(t, []string{"младше 18", "18-25", "26-40", "41-55", "старше 55"}, reply.options)
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
	require.Equal(t, fmt.Sprintf(welcomeHeader, 568)+"\n\nВы гражданин Республики Беларусь?", reply.text)
	require.Equal(t, []string{"Да", "Нет"}, reply.options)

	text := newTextCallback(t, userID, "Привет")
	require.Equal(t, text.User.ID, userID)
	reply, err = generateReplyFor(p, s, text)
	require.NoError(t, err)
	require.Equal(t, "Пожалуйста, выберите один из предложенных вариантов ответа.\n\nВы гражданин Республики Беларусь?", reply.text)
	require.Equal(t, []string{"Да", "Нет"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Да"))
	require.NoError(t, err)
	require.Equal(t, "Укажите, пожалуйста, Ваш возраст", reply.text)
	require.Equal(t, []string{"младше 18", "18-25", "26-40", "41-55", "старше 55"}, reply.options)

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
	require.Equal(t, reply.text, fmt.Sprintf(welcomeHeader, 568)+"\n\nВы гражданин Республики Беларусь?")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста, выберите один из предложенных вариантов ответа.\n\nВы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"Да", "Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "нету"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста, выберите один из предложенных вариантов ответа.\n\nВы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"Да", "Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "ДА"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "младше 18"))
	require.NoError(t, err)
	require.Equal(t, "Вам должно быть 18 или больше.\n\nУкажите, пожалуйста, Ваш возраст", reply.text)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "41-55"))
	require.NoError(t, err)
	require.Equal(t, "К какому типу относится населенный пункт, в котором вы проживаете?", reply.text)
	require.Equal(t, reply.options, []string{"Областной центр", "Город или городской поселок", "Агрогородок / Село / Деревня", "Проживаю за пределами РБ"})
	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Агрогородок / село / Деревня"))
	require.NoError(t, err)

	require.Equal(t, "За кого Вы планируете проголосовать?", reply.text)
	require.Equal(t, []string{"Дмитриев", "Тихановская", "Лукашенко", "Канопацкая", "Черечень", "Против всех", "Затрудняюсь ответить", "Не пойду голосовать"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Дмитриев"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш пол")
	require.Equal(t, reply.options, []string{"Мужской", "Женский"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "мужской"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Выберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск")
	require.Equal(t, reply.options, []string{"Брестская", "Витебская", "Гомельская", "Гродненская", "Минская", "Могилевская", "Минск", "Проживаю за пределами РБ"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "витебская"))
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equal(t, "Ваш уровень образования?", reply.text)
	require.Equal(t, []string{"Базовое / Среднее общее (школа)", "Профессионально-техническое", "Среднее специальное", "Высшее", "Другое"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Среднее специальное"))
	require.NoError(t, err)
	require.Equal(t, "Укажите, пожалуйста, общий совокупный месячный доход вашей семьи (включая пенсии, стипендии, надбавки и прочее)", reply.text)
	require.Equal(t, reply.options, []string{"До 500 бел. руб.", "500 - 1000 бел. руб.", "1000 - 2000 бел. руб.", "Выше 2000 бел.руб.", "Не хочу отвечать на этот вопрос"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "До 500 бел. руб."))
	require.NoError(t, err)
	require.Equal(t, "Спасибо за участие в нашем опросе!\nМы опубликуем результаты до 4 августа.\nНас уже 569 человек!", reply.text)

	user, err := s.fromPersisted(userID)
	require.NoError(t, err)
	require.Equal(t, user.ID, userID)
	require.Equal(t, user.Properties["age"], "41-55")
	require.Equal(t, 8, user.Level)
	require.Equal(t, user.Candidate, "Дмитриев")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Спасибо за участие в нашем опросе!\nМы опубликуем результаты до 4 августа.\nНас уже 569 человек!")

	reply, err = generateReplyFor(p, s, newUnsubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Nil(t, reply)

	reply, err = generateReplyFor(p, s, newSubscribeCallback(t, userID))
	require.NoError(t, err)
	require.Equal(t, "Добрый день!\nВы уже приняли участие в Народном опросе. Спасибо, ваш голос учтен!\nМы опубликуем результаты до 4 августа.\nНас уже 569 человек!", reply.text)
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
	require.Equal(t, fmt.Sprintf(welcomeHeader, 568)+"\n\nВы гражданин Республики Беларусь?", reply.text)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Привет"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Пожалуйста, выберите один из предложенных вариантов ответа.\n\nВы гражданин Республики Беларусь?")
	require.Equal(t, reply.options, []string{"Да", "Нет"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Да"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "младше 18"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Вам должно быть 18 или больше.\n\nУкажите, пожалуйста, Ваш возраст")

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "41-55"))
	require.NoError(t, err)

	require.Equal(t, "К какому типу относится населенный пункт, в котором вы проживаете?", reply.text)

	user, err := s.fromPersisted(userID)
	require.NoError(t, err)

	require.Equal(t, user.ID, userID)
	fmt.Println(user)
	require.Equal(t, user.Properties["age"], "41-55")
	require.Equal(t, 2, user.Level)

	seenCallback := newSeenCallback(t, userID)
	require.Equal(t, seenCallback.User.ID, userID)
	reply, err = generateReplyFor(p, s, seenCallback)
	require.NoError(t, err)
	require.Nil(t, reply)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Агрогородок / Село / Деревня"))
	require.NoError(t, err)
	require.Equal(t, "За кого Вы планируете проголосовать?", reply.text)
	require.Equal(t, []string{"Дмитриев", "Тихановская", "Лукашенко", "Канопацкая", "Черечень", "Против всех", "Затрудняюсь ответить", "Не пойду голосовать"}, reply.options)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Лукашенко"))
	require.NoError(t, err)
	require.Equal(t, reply.text, "Укажите, пожалуйста, Ваш пол")
	require.Equal(t, reply.options, []string{"Мужской", "Женский"})

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Мужской"))
	require.NoError(t, err)

	user, err = s.fromPersisted(userID)
	require.NoError(t, err)

	require.Equal(t, user.ID, userID)
	require.Equal(t, user.Properties["age"], "41-55")
	require.Equal(t, 5, user.Level)
	require.Equal(t, "Лукашенко", user.Candidate)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, "Пожалуйста, выберите один из предложенных вариантов ответа.\n\nВыберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск", reply.text)

	reply, err = generateReplyFor(p, s, newTextCallback(t, userID, "Передумал"))
	require.NoError(t, err)
	require.Equal(t, "Пожалуйста, выберите один из предложенных вариантов ответа.\n\nВыберите область, в которой Вы проживаете. Если Вы проживаете в Минске, выберите Минск", reply.text)

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
