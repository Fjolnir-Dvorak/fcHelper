using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.ComponentModel;
using System.Linq;
using System.Runtime.CompilerServices;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Threading;

namespace fcHelper
{
    public class LogLanguageEntry : INotifyPropertyChanged
    {

        private string language;

        public string Language
        {
            get => this.language;
            set {
                this.language = value;
                RaiseChange();
            }
        }
        private int min;
        public string Min
        {
            get => this.min.ToString();
            set
            {
                this.min = int.Parse(value);
                RaiseChange();
            }
        }
        private int max;
        public string Max
        {
            get => this.max.ToString();
            set
            {
                this.max = int.Parse(value);
                RaiseChange();
            }
        }
        private int value;
        public string Value
        {
            get => this.value.ToString();
            set
            {
                this.value = int.Parse(value);
                RaiseChange();
            }
        }

        private string message;

        public string Message
        {
            get => this.message;
            set {
                this.message = value;
                RaiseChange();
            }
        }

        public event PropertyChangedEventHandler PropertyChanged;

        private void RaiseChange([CallerMemberName] string caller = "")
        {
            PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(caller));
        }



        public LogLanguageEntry(string lang, int st, int en, int vl, string msg)
        {
            this.Language = lang;
            this.Min = st.ToString();
            this.Max = en.ToString();
            this.Value = vl.ToString();
            this.Message = msg;
        }

        [MethodImpl(MethodImplOptions.Synchronized)]
        public void increaseProcess()
        {
            this.Value = (this.value + 1).ToString();
            //Dispatcher.CurrentDispatcher.Invoke(new Action(() => { }), DispatcherPriority.ContextIdle, null);
        }

        public static ObservableCollection<LogLanguageEntry> GenTestCollection()
        {
            var collection = new ObservableCollection<LogLanguageEntry>();
            //collection.Add(new LogLanguageEntry("ASDF", 0, 100, 75, "Ich bin die erste Testnachricht"));
            //collection.Add(new LogLanguageEntry("QWERTZ", 0, 100, 75, "Ich bin die zweite Testnachricht"));
            return collection;
        }
    }
}
