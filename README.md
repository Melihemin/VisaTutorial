# Randevu Uygunluk Takip Sistemi

## Teknik Proje Dokümantasyonu

### Proje Özeti

Proje kapsamında farklı randevu platformlarının çalışma yapıları incelenmiş ve özellikle iDATA ile BLS International benzeri sistemlerde kontrollü testler gerçekleştirilmiştir.

Temel olarak çözmeye çalıştığım problem şuydu: Bir kullanıcı platform üzerinde gerekli doğrulama ve giriş işlemlerini tamamladıktan sonra, mevcut oturum bağlamını kullanarak randevu uygunluk durumunu düzenli ve kontrollü şekilde takip edebilen bir yapı oluşturmak.

Bu süreçte özellikle CAPTCHA, oturum süreleri, istek limitleri, kullanıcı bazlı kısıtlamalar ve ağ bağlantısının sürekliliği gibi gerçek sistemlerde karşılaşılan problemleri deneyimleme fırsatım oldu.

---

# 1. Teknoloji Seçimi

Projenin ana programlama dili olarak **Go** kullanıldı.

Go'yu seçmemdeki temel nedenler:

* Yüksek performanslı HTTP istemcileri oluşturabilmek
* Goroutine yapısı sayesinde eş zamanlı işlemleri kolay yönetebilmek
* Worker tabanlı mimariye uygun olması
* Ağ ve bağlantı işlemlerinde güçlü standart kütüphanelere sahip olması
* Uzun süre çalışan backend servisleri için sade bir yapı sunması

Proje boyunca mümkün olduğunca farklı sorumlulukları birbirinden ayırmaya çalıştım.

---

# 2. Genel Sistem Mimarisi

Sistem genel olarak aşağıdaki bileşenlerden oluşmaktadır:

```text
┌─────────────────┐
│    Kullanıcı     │
│  Manuel Giriş   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Proxy Worker   │
│   Sticky IP     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Session Manager │
│ Oturum Verileri │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Request Worker  │
│ Slot Sorgulama  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Response Parser │
│ Sonuç Analizi   │
└─────────────────┘
```

Her bileşenin farklı bir sorumluluğu bulunmaktadır. Böylece ağ bağlantısı, oturum yönetimi ve API sorguları birbirinden mümkün olduğunca bağımsız şekilde geliştirilebilmektedir.

---

# 3. Proxy ve Sticky IP Yapısı

Uygulama başlatıldığında ilk olarak proxy bağlantısı yapılandırılmaktadır.

Test mimarisinde, hedef sistemlerin bulunduğu coğrafi bölge ve kullanım senaryosuna uygun şekilde yapılandırılmış proxy altyapıları kullanılmıştır.

Burada özellikle **Sticky IP** yaklaşımı üzerinde çalışılmıştır.

Sticky IP kullanımı sayesinde aynı oturum boyunca bağlantının mümkün olduğunca aynı IP üzerinden devam etmesi amaçlanmaktadır.

Bu yaklaşımın temel amacı, oturum sırasında ağ bağlamının sürekli değişmesini önlemektir.

Genel bağlantı yapısı şu şekildedir:

```text
Go Application
       │
       ▼
Proxy Worker
       │
       ▼
Sticky IP Connection
       │
       ▼
Target Platform
```

Proxy katmanının ayrı bir worker olarak tasarlanması, ağ yapılandırmasının uygulamanın diğer bileşenlerinden ayrılmasını sağlamaktadır.

---

# 4. Kullanıcı Girişi ve Oturum Yönetimi

Sistem, kullanıcının hesabına otomatik olarak giriş yapacak şekilde tasarlanmamıştır.

Kullanıcı öncelikle ilgili platformda gerekli giriş ve doğrulama işlemlerini manuel olarak tamamlamaktadır.

Başarılı bir oturum oluşturulduktan sonra uygulamanın ihtiyaç duyduğu oturum ve istemci bağlamı yönetilmektedir.

Bu kapsamda aşağıdaki bilgiler üzerinde çalışılmıştır:

```text
SESSION_COOKIE
AUTHORIZE_TOKEN
CLIENT_SOURCE_TOKEN
USER_AGENT
SEC_CH_UA
```

Bu verilerin yönetimi için Go ile ayrı bir session management bileşeni oluşturulmuştur.

Temel amaç, aktif oturumun ihtiyaç duyduğu bilgileri merkezi bir noktada yönetmek ve gerekli olduğunda request worker'larına aktarmaktır.

---

# 5. Oturum Sürekliliği

Proje sırasında karşılaştığım önemli problemlerden biri oturumların ne kadar süre geçerli kaldığıydı.

Bazı sistemlerde kullanıcı başarılı bir şekilde doğrulama yaptıktan sonra mevcut oturum bir süre aktif kalırken, belirli sayıda yenileme veya belirli davranışlardan sonra kullanıcı yeniden giriş ekranına yönlendirilebilmektedir.

Bu nedenle oturum yönetimi tarafında aşağıdaki senaryoların değerlendirilmesi planlanmıştır:

* Oturumun geçerlilik durumunun kontrol edilmesi
* Geçersiz oturumların tespit edilmesi
* API yanıtlarındaki yetkilendirme hatalarının izlenmesi
* Kullanıcıya yeniden doğrulama gerektiğinin bildirilmesi
* Gereksiz sayfa yenilemelerinin önlenmesi

Bu noktada yapılan önemli gözlemlerden biri, sürekli sayfa yenilemenin oturum davranışını olumsuz etkileyebilmesidir.

Bu nedenle mimarinin, mümkün olduğunda kullanıcı arayüzünü sürekli yenilemek yerine oturumun ve yetkilendirmenin durumunu kontrollü şekilde takip edecek biçimde geliştirilmesi hedeflenmiştir.

---

# 6. Request Worker ve Slot Sorgulama Yapısı

Kullanıcı oturumu oluşturulduktan sonra sistemde ayrı bir **Request Worker** devreye girmektedir.

Bu worker'ın temel görevi, randevu uygunluk bilgisini sağlayan API uç noktalarından gelen sonuçları belirli kurallar çerçevesinde kontrol etmektir.

Genel çalışma akışı:

```text
Worker Başlatılır
       │
       ▼
Oturum Durumu Kontrol Edilir
       │
       ▼
İstek Parametreleri Hazırlanır
       │
       ▼
Bekleme Politikası Uygulanır
       │
       ▼
API Sorgusu Yapılır
       │
       ▼
JSON Yanıtı İşlenir
       │
       ├── Boş Slot Bulundu
       │
       └── Sonraki Sorgu Döngüsü
```

Sorguların sürekli ve kontrolsüz şekilde gönderilmesi yerine, belirli bekleme süreleri ve hata durumlarını dikkate alan bir yapı üzerinde çalışılmıştır.

---

# 7. API Yanıt Yapısı

BLS testleri sırasında aşağıdaki yapıya benzer JSON yanıtları elde edilmiştir:

```json
{
    "earliestDate": "08/21/2026 00:00:00",
    "earliestSlotLists": [
        {
            "applicant": "1,2",
            "date": "08/21/2026 00:00:00"
        },
        {
            "applicant": "3",
            "date": "08/24/2026 00:00:00"
        }
    ],
    "error": null
}
```

Bu yanıt yapısında sistemin değerlendirdiği temel bilgiler şunlardır:

* En erken uygun randevu tarihi
* Farklı başvuru grupları için uygunluk durumu
* İlgili tarihler
* Hata bilgisi

Go uygulaması gelen JSON verisini işleyerek uygun randevu tarihlerini tespit edecek şekilde yapılandırılmıştır.

---

# 8. Başarılı Test Sonucu

Test uygulaması sırasında uygun bir randevu bilgisinin başarıyla tespit edildiği gözlemlenmiştir.

Terminal çıktısı aşağıdaki şekildedir:

```text
[20:30:09] VACHOUS merkezi için takvim sorgulanıyor...

BOŞ SLOT BULUNDU! En Erken Tarih: 08/21/2026 00:00:00
  -> Başvuru Kombinasyonu: 1 | Tarih: 08/21/2026 00:00:00

[20:30:55] VACHOUS merkezi için takvim sorgulanıyor...

BOŞ SLOT BULUNDU! En Erken Tarih: 08/21/2026 00:00:00
  -> Başvuru Kombinasyonu: 1 | Tarih: 08/21/2026 00:00:00
```

Bu test, oluşturulan worker yapısının API'den gelen yanıtı işleyerek uygun tarih bilgisini başarılı bir şekilde ayrıştırabildiğini göstermiştir.

---

# 9. Karşılaşılan Rate Limiting Problemleri

Proje sırasında farklı platformların kullanıcı ve istek bazlı güvenlik mekanizmalarıyla karşılaşıldı.

Örneğin bazı testlerde sistem tarafından çok fazla isteğin kısa süre içerisinde gönderildiği tespit edilerek geçici erişim kısıtlamaları uygulanmıştır.

Örnek hata:

```text
Permission Issue (429201)

It seems like you have made multiple requests within a short period
which exceeds the defined thresholds.
```

Bu durum, gerçek dünya sistemlerinde **rate limiting** ve **request throttling** mekanizmalarının ne kadar önemli olduğunu gösterdi.

Bu gözlem sonucunda sistem mimarisinde aşağıdaki geliştirmelerin gerekli olduğu değerlendirildi:

* İstek sıklığının merkezi olarak yönetilmesi
* Hata kodlarının takip edilmesi
* Geçici kısıtlamalarda sorguların durdurulması
* Kontrollü cooldown mekanizmaları
* Gereksiz tekrar denemelerinin engellenmesi
* Kullanıcı bazlı hata durumlarının ayrıştırılması

---

# 10. Oturum ve Eski Verilerin Geçerliliği

Test sürecinde karşılaşılan bir diğer önemli durum, daha önce oluşturulmuş oturum bilgilerinin tekrar kullanılmasıydı.

Eski veya geçerliliğini kaybetmiş oturum bilgilerinin kullanılması sonucunda platform tarafından kullanıcı erişimine yönelik güvenlik kısıtlamaları uygulanabildiği gözlemlendi.

Bu nedenle sistem içerisinde aşağıdaki mekanizmaların bulunması önemlidir:

```text
Session Oluşturuldu
        │
        ▼
Session Geçerlilik Kontrolü
        │
        ├── Geçerli → İşleme Devam Et
        │
        ▼
Geçersiz / Süresi Dolmuş
        │
        ▼
Worker'ı Durdur
        │
        ▼
Kullanıcıdan Yeniden Doğrulama İste
```

Bu yaklaşım sayesinde geçersiz bir oturumla sürekli istek göndermek yerine sistemin kontrollü şekilde durması hedeflenmektedir.

---

# 11. HTTP ve TLS Bağlantı Yapısı

Proje sırasında HTTP istemci davranışı ve TLS bağlantı parametreleri üzerinde de çalışılmıştır.

Bu bölümde temel olarak aşağıdaki konular değerlendirilmiştir:

* HTTPS bağlantı yönetimi
* TLS yapılandırması
* HTTP timeout değerleri
* Request header yönetimi
* Bağlantı hatalarının işlenmesi
* Response parsing

Mevcut test ortamında kullanılan şifreleme yapılandırması proje gereksinimlerine uygun olarak değerlendirilmiş ve ilgili bağlantı profili **120** seviyesinde tutulmuştur.

Buradaki temel amaç, bağlantı katmanını güvenli ve tutarlı bir şekilde yönetmektir.

---

# 12. Test Sırasında Kullanılan Örnek Parametre Yapısı

Randevu sorgulama işlemlerinde kullanılan parametreler hedef sisteme göre değişiklik gösterebilmektedir.

Örnek bir test isteği aşağıdaki yapıya sahiptir:

```json
{
    "countryCode": "usa",
    "missionCode": "che",
    "vacCode": "NYC",
    "visaCategoryCode": "SUSVFF",
    "roleName": "Individual",
    "payCode": ""
}
```

---

# 13. Worker Tabanlı Mimari

Projeyi geliştirirken tüm işlemleri tek bir uygulama döngüsüne koymak yerine farklı görevleri birbirinden ayırmaya çalıştım.

Mevcut mimaride temel olarak şu bileşenler bulunmaktadır:

### Proxy Worker

Ağ bağlantısının yönetiminden sorumludur.

### Session Manager

Kullanıcı oturumuna ait gerekli bilgileri yönetir ve oturum durumunu takip eder.

### Request Worker

Belirlenen sorgu politikasına göre API iletişimini gerçekleştirir.

### Response Parser

API'den gelen JSON yanıtlarını işleyerek kullanılabilir sonuçlara dönüştürür.

Bu yapı aşağıdaki şekilde özetlenebilir:

```text
┌─────────────────────┐
│    Proxy Worker     │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   Session Manager   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   Request Worker    │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   Response Parser   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Notification Layer  │
└─────────────────────┘
```

---


# Sonuç

Bu proje, benim için Go kullanarak gerçek dünya problemlerine yönelik worker tabanlı bir backend mimarisi geliştirme çalışması oldu.

Proje sürecinde sadece başarılı API istekleri üzerinde değil; aynı zamanda **oturumların sona ermesi, eski oturum verilerinin geçersizleşmesi, rate limiting, geçici erişim kısıtlamaları ve ağ bağlantısının sürekliliği** gibi başarısızlık senaryoları üzerinde de çalıştım.

Bu deneyimler sonucunda, başarılı bir backend sisteminin yalnızca doğru isteği göndermekten ibaret olmadığını; oturum yaşam döngüsünü, hata durumlarını, kaynak kullanımını ve sistemin güvenlik sınırlarını da doğru şekilde yönetmesi gerektiğini daha iyi anlamış oldum.

Bu proje benim için özellikle **Go, concurrency, network programming, session management ve API integration** alanlarında önemli bir uygulamalı deneyim sağlamıştır.

