public class Test{
	public static void main(String[] args){
		String a = "Hello, MLee";
		String b = "Hello, MLee";
		String c = new String("Hello, MLee");
		String d = new String("Hello, MLee");
		System.out.println("a == b: " +  a==b);
		System.out.println("a == c: " +  a==c);
		System.out.println("c == d: " +  c==d);
		System.out.println("a.equals(b): " +  a.equals(b));
		System.out.println("a.equals(c): " +  a.equals(c));
		System.out.println("c.equals(d): " +  c.equals(d));
	}
}
