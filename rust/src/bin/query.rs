use qnip::query::QueryDb;

fn main() {
    let args = std::env::args().collect::<Vec<String>>();
    println!("args: {:?}", args);
    if args.len() < 3 {
        println!("Usage: {} <ipdb> <ip>", args[0]);
        return;
    }
    let q = QueryDb::from_file(&args[1]);
    match q {
        Ok(qdb) => {
            let ip = &args[2];
            let info = qdb.query_str(ip);
            match info {
                Ok(info) => {
                    println!("query ip: {}, info: {}", ip, info);
                },
                Err(e) => {
                    println!("query ip: {}, error: {:?}", ip, e);
                },
            }
        },
        Err(e) => {
            println!("open ipdb error: {:?}", e);
        },
    }
}