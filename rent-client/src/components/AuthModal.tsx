
import React, { useState, useRef, useEffect} from 'react';
import { useUserStore } from '../stores/useUserStore';
import { Modal, Button, Input } from './UI';
import googleLoginIcon from '../assets/google_icon.svg';
import { postToBroker } from "../api/broker";

type AuthStep = 'SIGN_IN' | 'SIGN_UP' | 'FORGOT_PASSWORD' | 'VERIFY_CODE' | 'NEW_PASSWORD';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
  onLogin: (user: { first_name: string; last_name: string; email: string }) => void;
}

export const AuthModal: React.FC<AuthModalProps> = ({ isOpen, onClose, onLogin }) => {
  const [step, setStep] = useState<AuthStep>('SIGN_IN');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [first_name, setFirstName] = useState('');
  const [last_name, setLastName] = useState('');
  const [isSignUpFlow, setIsSignUpFlow] = useState(false);
  const [error, setError] = useState("");
  const [otp, setOtp] = useState(['', '', '', '', '', '']);
  const otpRefs = useRef<(HTMLInputElement | null)[]>([]);
  const setUser = useUserStore(state => state.setUser);
	const user = useUserStore(state => state.user);

  // Focus management for OTP
  useEffect(() => {
    if (step === 'VERIFY_CODE') {
      // Small delay to ensure modal transition finished
      setTimeout(() => otpRefs.current[0]?.focus(), 150);
    }
  }, [step]);

  const handleOtpChange = (value: string, index: number) => {
    if (!/^\d*$/.test(value)) return;
    
    const newOtp = [...otp];
    newOtp[index] = value.slice(-1);
    setOtp(newOtp);

    // Auto-advance
    if (value && index < 5) {
      otpRefs.current[index + 1]?.focus();
    }
  };

  const handleOtpKeyDown = (e: React.KeyboardEvent, index: number) => {
    if (e.key === 'Backspace' && !otp[index] && index > 0) {
      otpRefs.current[index - 1]?.focus();
    }
  };


  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (step === 'SIGN_IN') {
		console.log("User before", user);
	  const payload = { action: "auth", auth: { email, password } };
		const data = await postToBroker(payload);
		if (data.error) {
			setError(`${data.message}`);
			console.error(`Auth error: ${data.message}`);
		}else{
			console.log("Login successful:", data.data);
			onLogin(data.data);
			setUser({
					...data.data
				});

			console.log("User set in store:", useUserStore.getState().user); 
			//getState只适合非ui组件调试使用
			setError("");
			onClose();
		}
		
    } else if (step === 'SIGN_UP') {
			const payload = {
				action: "verify",
					verify: {
						email,
					},
				};
			console.log("Signup payload:", payload);
			const data = await postToBroker(payload);
		if (data.error) {
			console.error(`Sending code error: ${data.message}`);
			setError(`${data.message}`);
		} else{
			// Now requires verification
      		setIsSignUpFlow(true);
      		setStep('VERIFY_CODE');
		} 
    } else if (step === 'FORGOT_PASSWORD') {
		console.log("Forgot password for:", email);
      setIsSignUpFlow(false);
      setStep('VERIFY_CODE');
    } else if (step === 'VERIFY_CODE') {
		const payload = {
		action: "register",
		  register: {
		  email,
		  password,
		  first_name,
		  last_name,
		},
    	};
		console.log("verify:", payload)
		console.log(otp)
      if (isSignUpFlow) {
        onLogin({
          first_name: 'u1',
          last_name: 'Voss',
          email: email || 'julian@example.com'
        });
        onClose();
      } else {
        setStep('NEW_PASSWORD');
      }
    } else if (step === 'NEW_PASSWORD') {
      setStep('SIGN_IN');
    }
  };

  const renderTitle = () => {
    switch (step) {
      case 'SIGN_IN': return "SIGN IN";
      case 'SIGN_UP': return "JOIN DWELL";
      case 'FORGOT_PASSWORD': return "RESET ACCESS";
      case 'VERIFY_CODE': return "VERIFY IDENTITY";
      case 'NEW_PASSWORD': return "NEW PASSWORD";
    }
  };

  return (
    <Modal 
      isOpen={isOpen} 
      onClose={() => {
        onClose();
        setStep('SIGN_IN');
        setOtp(['', '', '', '', '', '']);
      }} 
      title={renderTitle()}
    >
      <div>
        <form onSubmit={handleSubmit} className="space-y-6">
	  	  {step === 'SIGN_UP' && <Input placeholder="FIRST NAME" type="text" value={first_name}
              onChange={(e) => setFirstName(e.target.value)} required />}

          {step === 'SIGN_UP' && <Input placeholder="LAST NAME" type="text" value={last_name}
              onChange={(e) => setLastName(e.target.value)} required />}
          
          {(step === 'SIGN_IN' || step === 'SIGN_UP' || step === 'FORGOT_PASSWORD') && (
            <Input 
              placeholder="EMAIL ADDRESS" 
              type="email" 
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required 
            />
          )}

          {step === 'VERIFY_CODE' && (
            <div className="space-y-8 py-4">
              <div className="text-center space-y-2">
                <p className="text-[12px] text-[#7e918b] uppercase tracking-[0.3em] leading-relaxed">
                  Verification code sent to
                </p>
                <p className="text-[#4a586e] font-bold text-xs uppercase tracking-tighter">
                  {email || 'your email address'}
                </p>
              </div>
              
              <div className="flex justify-between gap-2 max-w-xs mx-auto">
                {otp.map((digit, idx) => (
                  <input
                    key={idx}
                    ref={el => { otpRefs.current[idx] = el; }}
                    type="text"
                    maxLength={1}
                    value={digit}
                    onChange={(e) => handleOtpChange(e.target.value, idx)}
                    onKeyDown={(e) => handleOtpKeyDown(e, idx)}
                    className="w-10 h-14 md:w-12 md:h-16 bg-transparent border border-[#4a586e]/20 text-center text-2xl font-bold focus:border-[#4a586e] focus:outline-none transition-all text-[#4a586e] rounded-none"
                  />
                ))}
              </div>
              
              <p className="text-center text-[9px] text-[#7e918b] uppercase tracking-[0.2em]">
                Didn't receive code? <button type="button" className="text-[#4a586e] font-bold hover:underline">Resend</button>
              </p>
            </div>
          )}

          {step === 'NEW_PASSWORD' && (
            <>
              <Input placeholder="NEW PASSWORD" type="password" required />
              <Input placeholder="CONFIRM NEW PASSWORD" type="password" required />
            </>
          )}

          {(step === 'SIGN_IN' || step === 'SIGN_UP') && (
            <Input placeholder="PASSWORD" type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)} required />
          )}
          
          <Button type="submit" className="w-full mt-4">
            {step === 'SIGN_IN' && "CONTINUE"}
            {step === 'SIGN_UP' && "VERIFY EMAIL"}
            {step === 'FORGOT_PASSWORD' && "SEND CODE"}
            {step === 'VERIFY_CODE' && (isSignUpFlow ? "FINISH REGISTRATION" : "VERIFY & CONTINUE")}
            {step === 'NEW_PASSWORD' && "UPDATE PASSWORD"}
          </Button>
        </form>

		{error !== "" && <p className='text-[#f47979] my-0.5'>{error}</p>}

        {step === 'SIGN_IN' && (
          <div className="text-right">
            <button 
              onClick={() => setStep('FORGOT_PASSWORD')}
              className="text-[10px] uppercase tracking-widest text-[#7e918b] hover:text-[#4a586e] transition-colors font-bold border-b border-transparent hover:border-[#4a586e]"
            >
              Forgot Password?
            </button>
          </div>
        )}

        {(step === 'SIGN_IN' || step === 'SIGN_UP') && (
          <>
            <div className="relative py-4">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-[#4a586e]/10"></div>
              </div>
              <div className="relative flex justify-center text-xs uppercase tracking-widest">
                <span className="bg-[#f3e9d2] px-4 text-[#7e918b]">or</span>
              </div>
            </div>

            <button 
              type="button"
              onClick={handleSubmit}
              className="w-full flex items-center justify-center gap-4 border border-[#4a586e] p-4 hover:bg-[#4a586e] hover:text-[#f3e9d2] transition-colors text-[#4a586e]"
            >
              <img src={googleLoginIcon} alt="Google login icon" />
              <span className="text-[10px] font-bold uppercase tracking-widest">Continue with Google</span>
            </button>
          </>
        )}

        <p className="text-center text-[10px] uppercase tracking-widest text-[#7e918b]">
          {step === 'SIGN_IN' ? "New to Dwell?" : "Return to access?"}
          <button 
            type="button"
            onClick={() => setStep(step === 'SIGN_IN' ? 'SIGN_UP' : 'SIGN_IN')}
            className="ml-2 text-[#4a586e] hover:underline font-bold"
          >
            {step === 'SIGN_IN' ? "Sign up" : "Sign in"}
          </button>
        </p>

        {(step === 'FORGOT_PASSWORD' || step === 'VERIFY_CODE' || step === 'NEW_PASSWORD') && (
          <button 
            type="button"
            onClick={() => setStep('SIGN_IN')}
            className="w-full text-center text-[8px] uppercase tracking-widest text-[#7e918b] hover:text-[#4a586e] font-bold transition-colors"
          >
            ← Back to Sign In
          </button>
        )}
      </div>
    </Modal>
  );
};


